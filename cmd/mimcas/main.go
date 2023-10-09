package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

//// lru cache data structure

type Cache struct {
	items       map[string]*Node
	lruList     *list.List
	memory 	    int
	maxmemory   int
	emptyCacheSizeBytes int
	emptyItemSizeBytes  int
}

type Node struct {
	mutex sync.RWMutex
	key   string
	value string
	lruElem *list.Element
}

//// inserts channel and handler function
var inserts = make(chan *Node, 100) // TODO: parameterize channel sizes. Must be different depending on the load.
func insert(toInsert *Node) { inserts <- toInsert }

func (c *Cache) insertsHandler() {
	for {
		toInsert := <-inserts
		if node, ok := c.items[toInsert.key]; ok { // to prevent inserting twice for a key
			node.mutex.Lock()
			node.value = toInsert.value
			node.mutex.Unlock()
			if (0 < c.maxmemory){
				markAsUsed(node)
			}
		} else {
			insertLru(toInsert)
			c.items[toInsert.key] = toInsert
		}
	}
}

//// lruOperationsHandler goroutine and stuff related to it
type lruOperationsChanStruct struct {
	op int	// 0 mark as used
			// 1 insert
			// 2 remove
	el *list.Element 	// for op 2
	n *Node				// for op 0 and 1
}
var lruOperations = make(chan lruOperationsChanStruct, 100)
func markAsUsed(node *Node) { lruOperations <-lruOperationsChanStruct{n: node, op: 0} }
func insertLru(node *Node) { lruOperations <- lruOperationsChanStruct{n: node, op: 1} }
func removeLru(elem *list.Element) { lruOperations <- lruOperationsChanStruct{el: elem, op: 2} }

func (c *Cache) lruOperationsHandler() {
	for {
		lruOperation := <-lruOperations
		switch lruOperation.op {
			case 0: // mark as used
				c.lruList.MoveToFront(lruOperation.n.lruElem)
				continue
			case 1: // insert
				elem := c.lruList.PushFront(lruOperation.n)
				elem.Value.(*Node).lruElem = elem
				continue
			//case 2: // delete
			default:
				fmt.Printf("Error: invalid lru operation code")
		}
	}
}

// function for measuring memory usage by the LRU list
var memoryDeltas = make(chan int, 100)
func memoryDelta(delta int) { memoryDeltas <- delta }
func (c *Cache) memoryHandler() {
	c.emptyItemSizeBytes = int(unsafe.Sizeof(c.lruList.Front()) + unsafe.Sizeof(c.lruList.Front().Value.(*Node)) + unsafe.Sizeof(c.lruList.Front().Value.(*Node).mutex))
	c.emptyCacheSizeBytes = int(unsafe.Sizeof(*c) + unsafe.Sizeof(c.lruList))
	c.memory = c.emptyCacheSizeBytes

	for {
		delta := <-memoryDeltas
		c.memory += delta
		if (c.maxmemory < c.memory) {
			evict(c, c.memory - c.maxmemory)
		}
		//fmt.Println(c.memory) //useful for debugging
	}
}

func evict(c *Cache, bytesToEvict int) {
	bytesEvicted := 0
	for bytesEvicted < bytesToEvict{
		backElem := c.lruList.Back()
		c.lruList.Remove(backElem)
		delete(c.items, backElem.Value.(*Node).key)
		//fmt.Println("deleted: " + backElem.Value.(*Node).key)
		bytesEvicted += c.emptyItemSizeBytes + len(backElem.Value.(*Node).key) + len(backElem.Value.(*Node).value)
	}
	memoryDelta(0 - bytesEvicted)
}

//// commands

func (c *Cache) set(params []string) string {
	// This function is "if" hell. Refactoring it would be good
	response := ""
	if len(params) == 3 {
		itemSizeBytes := 0
		if node, ok := c.items[params[1]]; ok {
			delta := 0
			node.mutex.Lock()
			if (0 < c.maxmemory){
				itemSizeBytes = c.emptyItemSizeBytes + len(node.key) + len(params[2])
				delta = len(params[2]) - len(node.value)
			}
			if (c.maxmemory < itemSizeBytes + c.emptyCacheSizeBytes && 0 < c.maxmemory){
				response = "ERR item too big, increase maxmemory parameter.\n"
			} else {
				node.value = params[2]
				if (0 < c.maxmemory){
					memoryDelta(delta)
					markAsUsed(node)
				}
				response = "OK\n"
			}
			node.mutex.Unlock()
		} else {
			node := Node{key: params[1], value: params[2]}
			if (0 < c.maxmemory){
				itemSizeBytes = c.emptyItemSizeBytes + len(node.key) + len(node.value)
			}
			if (c.maxmemory < itemSizeBytes + c.emptyCacheSizeBytes && 0 < c.maxmemory){
				response = "ERR item too big, increase maxmemory parameter\n"
			} else {
				insert(&node)
				if (0 < c.maxmemory){
					memoryDelta(itemSizeBytes)
				}
				response = "OK\n"
			}
		}
	} else {
		response = "ERR syntax error\n"
	}
	return response
}

func (c *Cache) get(params []string) string {
	response := ""
	if len(params) == 2 {
		if node, ok := c.items[params[1]]; ok {
			node.mutex.RLock()
			value := node.value
			node.mutex.RUnlock()
			// mark as read could be non blocking by using a buffered channel,
			// but as a downside there is risk to not mark as used some used data
			// if the channel fills up too quickly
			if (0 < c.maxmemory){
				markAsUsed(node)
			}
			if value == "" {
				response = "(nil)\n"
			} else {
				response = value + "\n"
			}
		} else {
			response = "(nil)\n"
		}
	} else {
		response = "ERR syntax error\n"
	}
	return response
}

func (c *Cache) mget(params []string) string {
	response := ""
	if 2 <= len(params) {
		for _, key := range params[1:] {
			if node, ok := c.items[key]; ok {
				node.mutex.RLock()
				value := node.value
				node.mutex.RUnlock()
				if (0 < c.maxmemory){
					markAsUsed(node)
				}
				if value == "" {
					response = response + "(nil)\n"
				} else {
					response = response + value + "\n"
				}
			} else {
				response = response + "(nil)\n"
			}
		}
	} else {
		response = "ERR syntax error\n"
	}
	return response
}

func (c *Cache) delete(params []string) string {
	response := ""
	if len(params) == 2 {
		if node, ok := c.items[params[1]]; ok {
			node.mutex.Lock()
			delete(c.items, node.key)
			if (0 < c.maxmemory){
				c.lruList.Remove(node.lruElem)
				itemSizeBytes := c.emptyItemSizeBytes + len(node.key) + len(node.key)
				memoryDelta(0 - itemSizeBytes)
			}
			response = "OK\n"
			node.mutex.Unlock()
		} else {
			response = "(nil)\n"
		}
	} else {
		response = "ERR syntax error\n"
	}
	return response
}

//// connection handling and main

func handleConnection(cache *Cache, conn net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		data := strings.TrimSpace(string(netData))
		params := strings.Split(data, " ")
		response := ""
		switch params[0] {
		case "set":
			response = cache.set(params)
		case "get":
			response = cache.get(params)
		case "mget":
			response = cache.mget(params)
		case "del":
			response = cache.delete(params)
		case "quit":
			break
		case "ping":
			response = "pong\n"
		default: 
			response = "ERR unknown command\n"
		}
		conn.Write([]byte(response))
	}
	conn.Close()
	fmt.Printf("Closed conn to %s\n", conn.RemoteAddr().String())
}

func main() {
	var port int
	var maxmemory int
	flag.IntVar(&port, "port", 20000, "Port to use for listening for incoming connections.")
	flag.IntVar(&maxmemory, "maxmemory", -1, "Maximum number of bytes available to use. Items will be evicted following LRU policy when that limit is crossed. By default there is no limit.")
	flag.Parse()

	var cache = Cache{items: make(map[string]*Node), lruList: list.New(), maxmemory: maxmemory}

	go cache.insertsHandler()
	go cache.lruOperationsHandler()
	if (0 < cache.maxmemory) {
		go cache.memoryHandler()
	}

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error listening")
		os.Exit(2)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection")
			os.Exit(2)
		}
		go handleConnection(&cache, conn)
	}
}

package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"unsafe"
	"time"
)

//// lru cache data structure

type Cache struct {
	items       map[string]*list.Element
	lruList     *list.List
	memory 	    int
	maxmemory   int
	memoryMutex sync.RWMutex
}

type node struct {
	mutex sync.RWMutex
	key   string
	value string
}

//// insertsHandler goroutine and stuff related to it
type insertsChanStruct struct {
	n   *node
	key string
}

var inserts = make(chan insertsChanStruct)
func insert(toInsert insertsChanStruct) { inserts <- toInsert }

func (c *Cache) insertsHandler() {
	for {
		toInsert := <-inserts
		if elem, ok := c.items[toInsert.key]; ok { // to prevent inserting twice for a key
			elem.Value.(*node).mutex.Lock()
			elem.Value.(*node).value = toInsert.n.value
			elem.Value.(*node).mutex.Unlock()
			markAsUsed(elem)
		} else {
			n := toInsert.n
			insertLru(n)
			c.items[toInsert.key] = elem
		}
	}
}

//// lruOperationsHandler goroutine and stuff related to it
type lruOperationsChanStruct struct {
	op int	// 0 mark as used
			// 1 insert
			// 2 remove
	el *list.Element 	// for op 0 and 2
	n *node				// for op 1
}
var lruOperations = make(chan lruOperationsChanStruct)
func markAsUsed(elem *list.Element) { lruOperations <-lruOperationsChanStruct{el: elem, op: 0} }
func insertLru(node *node) { lruOperations <- lruOperationsChanStruct{n: node, op: 1} }
func removeLru(elem *list.Element) { lruOperations <- lruOperationsChanStruct{el: elem, op: 2} }

func (c *Cache) lruOperationsHandler() {
	for {
		lruOperation := <-lruOperations
		if (lruOperation.op == 0) { // mark as used
			c.lruList.MoveToFront(lruOperation.el)
		} else if (lruOperation.op == 1) { // insert
			c.lruList.PushFront(lruOperation.n)
		} else if (lruOperation.op == 2) { // delete

		} else {
			fmt.Printf("Error: invalid lru operation code")
		}

	}
}

// function for measuring memory usage by the LRU list
func (c *Cache) evictionHandler() {
	// TODO
	// * WARNING The list node order can change while this is happening. This is an provisional solution.
	//		- either lock the list while iterating or keep track when adding, editing or removing nodes
	// * WARNING not checking for overflows. Unsafe package is called that way for a reason
	// * WARNING probably I'm missing something, need to check the total memory footprint of the program
	lruElementsSizeBytes := unsafe.Sizeof(c.lruList.Front()) + unsafe.Sizeof(c.lruList.Front().Value.(*node))
	emptyCacheSizeBytes := int(unsafe.Sizeof(*c)) + int(unsafe.Sizeof(c.lruList))
	c.memoryMutex.Lock()
	c.memory = emptyCacheSizeBytes
	c.memoryMutex.Unlock()
	for {
	    nodesSizeBytes := c.lruList.Len() * int(lruElementsSizeBytes)
	    for e := c.lruList.Front(); e != nil; e = e.Next() {
	    	nodesSizeBytes += int(unsafe.Sizeof(e.Value.(*node)))
	    	nodesSizeBytes += int(unsafe.Sizeof(e.Value.(*node).mutex))
	    	nodesSizeBytes += len(e.Value.(*node).key)
	    	nodesSizeBytes += len(e.Value.(*node).value)
	    }
	    totalCacheSizeBytes := emptyCacheSizeBytes + nodesSizeBytes
	    fmt.Printf("totalCacheSizeBytes: %d\n", totalCacheSizeBytes)

		fmt.Printf("Size of cache: %d bytes\n", unsafe.Sizeof(*c))
		time.Sleep(30 * time.Second) // TODO make this a config value
	}
}

//// commands

func (c *Cache) set(params []string) string {
	response := ""
	if len(params) == 3 {
		if elem, ok := c.items[params[1]]; ok {
			elem.Value.(*node).mutex.Lock()
			elem.Value.(*node).value = params[2]
			elem.Value.(*node).mutex.Unlock()
			markAsUsed(elem)
		} else {
			newNode := node{key: params[1], value: params[2]}
			toInsert := insertsChanStruct{n: &newNode, key: params[1]}
			// insert could be non blocking by using a buffered channel,
			// but as a downside there is risk to loose inserted data
			// if the channel fills up too quickly
			insert(toInsert)
		}
		response = "OK\n"
	} else {
		response = "ERR syntax error\n"
	}
	return response
}

func (c *Cache) get(params []string) string {
	response := ""
	if len(params) == 2 {
		if elem, ok := c.items[params[1]]; ok {
			elem.Value.(*node).mutex.RLock()
			value := elem.Value.(*node).value
			elem.Value.(*node).mutex.RUnlock()
			// mark as read could be non blocking by using a buffered channel,
			// but as a downside there is risk to not mark as used some used data
			// if the channel fills up too quickly
			markAsUsed(elem)
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
			if elem, ok := c.items[key]; ok {
				elem.Value.(*node).mutex.RLock()
				value := elem.Value.(*node).value
				elem.Value.(*node).mutex.RUnlock()
				markAsUsed(elem)
				if value == "" {
					response = response + "(nil)\n"
				} else {
					response = response + value + "\n"
				}
			} else { // TODO something wrong here (?)
				response = response + "(nil)\n"
			}
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
		if params[0] == "SET" || params[0] == "set" {
			response = cache.set(params)
		} else if params[0] == "GET" || params[0] == "get" {
			response = cache.get(params)
		} else if params[0] == "MGET" || params[0] == "mget" {
			response = cache.mget(params)
		} else if params[0] == "QUIT" || params[0] == "quit" {
			break
		} else if params[0] == "PING" || params[0] == "ping" {
			response = "PONG\n"
		} else {
			response = "ERR unknown command\n"
		}
		conn.Write([]byte(response))
	}
	conn.Close()
	fmt.Printf("Closed conn to %s\n", conn.RemoteAddr().String())
}

func main() {
	var port int
	var apiport int
	var maxmemory int
	flag.IntVar(&port, "port", 20000, "port to listen to")
	flag.IntVar(&apiport, "apiport", 8080, "port to listen to")
	flag.IntVar(&maxmemory, "maxmemory", 0, "Maximum number of bytes available to use")
	flag.Parse()

	var cache = Cache{items: make(map[string]*list.Element), lruList: list.New()}

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong\n")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
	go http.ListenAndServe(":"+strconv.Itoa(apiport), nil)

	go cache.insertsHandler()
	go cache.lruOperationsHandler()
	if (maxmemory > 0) {
		go cache.evictionHandler()
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
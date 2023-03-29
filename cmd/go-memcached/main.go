package main

import (
	"container/list"
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

//// lru cache data structure

type Cache struct {
	items map[string]*list.Element
	lruList *list.List
	// should channels be inside the cache struct or be global?
}

type node struct {
	mutex sync.RWMutex
	key   string
	value string
}

//// cacheHandler goroutine and stuff related to it
type insertsChanStruct struct {
	n *node
	key string
}

var inserts = make(chan insertsChanStruct)
func insert(toInsert insertsChanStruct) { inserts <- toInsert }

var useds = make(chan *list.Element)
func markAsUsed(elem *list.Element) { useds <- elem }

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
			elem := c.lruList.PushFront(n)
			c.items[toInsert.key] = elem
		}
	}
}

func (c *Cache) usedsHandler() {
	for {
		elem := <-useds
		c.lruList.MoveToFront(elem)
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
			} else { // TODO something wrong here
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
	flag.IntVar(&port, "port", 20000, "port to listen to")
	flag.IntVar(&apiport, "apiport", 8080, "port to listen to")
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
	go cache.usedsHandler()

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

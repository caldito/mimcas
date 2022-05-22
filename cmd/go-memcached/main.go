package main

import (
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

//// basic cache data structure

type Cache struct {
	m map[string]*node
}

type node struct {
	mutex   sync.RWMutex
	value string
}

//// inserter goroutine and other related stuff

var inserts = make(chan insertsChanStruct)

func insert(a insertsChanStruct) { inserts <- a }

type insertsChanStruct struct {
	n node
	key string
}

func (c *Cache) inserter() {
	for {
		select {
		case a := <-inserts:
			if _, ok := c.m[a.key]; ok {
				c.m[a.key].mutex.Lock()
				c.m[a.key].value = a.n.value
				c.m[a.key].mutex.Unlock()
			} else {
				n := a.n
				c.m[a.key] = &n
			}
		}
	}
}

//// commands

func (c *Cache) set(params []string) string {
	response := ""
	if len(params) == 3 {
		if _, ok := c.m[params[1]]; ok {
			c.m[params[1]].mutex.Lock()
			c.m[params[1]].value = params[2]
			c.m[params[1]].mutex.Unlock()
		} else {
			newNode := node{value: params[2]}
			a := insertsChanStruct{n: newNode, key: params[1]}
			// insert could be non blocking by using a buffered channel,
			// but as a downside there is risk to loose inserted data
			// if the channel fills up too quickly
			insert(a)
		}
		response = "OK\n"
	} else {
		response = "ERR syntax error\n"
	}
	return response
}

func (c *Cache) get(params []string) string {
	response := ""
	if 2 == len(params) {
		if _, ok := c.m[params[1]]; ok {
			c.m[params[1]].mutex.RLock()
			value := c.m[params[1]].value
			c.m[params[1]].mutex.RUnlock()
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
			if _, ok := c.m[key]; ok {
				c.m[key].mutex.RLock()
				value := c.m[key].value
				c.m[key].mutex.RUnlock()
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

	var cache = Cache{m: make(map[string]*node)}

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong\n")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
	go http.ListenAndServe(":"+strconv.Itoa(apiport), nil)

	go cache.inserter()

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

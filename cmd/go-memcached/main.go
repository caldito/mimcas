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

type Node struct {
	Mu   sync.RWMutex
	Data string
}

var m = make(map[string]*Node)

//// inserter goroutine and other related stuff

var inserts = make(chan chanInsert)

func Insert(a chanInsert) { inserts <- a }

type chanInsert struct {
	N Node
	Key string
}

func inserter() {
	for {
		select {
		case a := <-inserts:
			if _, ok := m[a.Key]; ok {
				m[a.Key].Mu.Lock()
				m[a.Key].Data = a.N.Data
				m[a.Key].Mu.Unlock()
			} else {
				n := a.N
				m[a.Key] = &n
			}
		}
	}
}

//// commands

func set(params []string) string {
	response := ""
	if len(params) == 3 {
		if _, ok := m[params[1]]; ok {
			m[params[1]].Mu.Lock()
			m[params[1]].Data = params[2]
			m[params[1]].Mu.Unlock()
		} else {
			n := Node{Data: params[2]}
			a := chanInsert{N: n, Key: params[1]}
			// insert could be non blocking by using a buffered channel,
			// but as a downside there is risk to loose inserted data
			// if the channel fills up too quickly
			Insert(a)
		}
		response = "OK\n"
	} else {
		response = "ERR syntax error\n"
	}
	return response
}

func get(params []string) string {
	response := ""
	if 2 == len(params) {
		if _, ok := m[params[1]]; ok {
			m[params[1]].Mu.RLock()
			value := m[params[1]].Data
			m[params[1]].Mu.RUnlock()
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

func mget(params []string) string {
	response := ""
	if 2 <= len(params) {
		for _, key := range params[1:] {
			if _, ok := m[params[1]]; ok {
				m[key].Mu.RLock()
				value := m[key].Data
				m[key].Mu.RUnlock()
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

//// connection handling and main

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		data := strings.TrimSpace(string(netData))
		params := strings.Split(data, " ")
		response := ""
		if params[0] == "SET" || params[0] == "set" {
			response = set(params)
		} else if params[0] == "GET" || params[0] == "get" {
			response = get(params)
		} else if params[0] == "MGET" || params[0] == "mget" {
			response = mget(params)
		} else if params[0] == "QUIT" || params[0] == "quit" {
			break
		} else if params[0] == "PING" || params[0] == "ping" {
			response = "PONG\n"
		} else {
			response = "ERR unknown command\n"
		}
		c.Write([]byte(response))
	}
	c.Close()
	fmt.Printf("Closed conn to %s\n", c.RemoteAddr().String())
}

func main() {
	var port int
	var apiport int
	flag.IntVar(&port, "port", 20000, "port to listen to")
	flag.IntVar(&apiport, "apiport", 8080, "port to listen to")
	flag.Parse()

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	go http.ListenAndServe(":"+strconv.Itoa(apiport), nil)

	go inserter()

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
		go handleConnection(conn)
	}
}

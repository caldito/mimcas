package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"flag"
	"strconv"
	"net"
	"fmt"
	"os"
	"bufio"
	"strings"
)

var m = make(map[string]string)

func getKey(w http.ResponseWriter, r *http.Request) {
	keyParam := chi.URLParam(r, "key")
	value := m[keyParam]
	if value == "" {
		w.WriteHeader(404)
	} else {
		w.Write([]byte(value))
	}
}

func setKey(w http.ResponseWriter, r *http.Request) {
	keyParam := chi.URLParam(r, "key")
	valueParam := chi.URLParam(r, "val")
	m[keyParam] = valueParam
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		data := strings.TrimSpace(string(netData))
		params := strings.Split(data," ")
		response := ""
		if params[0] == "SET" {
			if len(params) == 3 {
				m[params[1]] = params[2]
				response = "OK\n"
			} else {
				response = "ERR syntax error\n"
			}
		} else if params[0] == "GET" {
			if 2 == len(params) {
				value := m[params[1]]
				if value == "" {
					response = "(nil)\n"
				} else {
					response = value + "\n"
				}
			} else {
				response = "ERR syntax error\n"
			}
		} else if params[0] == "MGET" {
			if 2 <= len(params) {
				for _, key := range params[1:] {
					value := m[key]
					if value == "" {
						response = response + "(nil)\n"
					} else {
						response = response + value + "\n"
					}
				}
			} else {
				response = "ERR syntax error\n"
			}
		} else if params[0] == "QUIT" {
			break
		} else if params[0] == "PING" {
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

	r := chi.NewRouter()
	r.Get("/keys/{key}", getKey)
	r.Put("/keys/{key}/{val}", setKey)
	r.Get("/health", getHealth)
	go http.ListenAndServe(":" + strconv.Itoa(apiport), r)

	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
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

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
)

var m = make(map[string]string)

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

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	go http.ListenAndServe(":"+strconv.Itoa(apiport), nil)

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

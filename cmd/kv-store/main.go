package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"flag"
	"strconv"
	// "time"
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
		if params[0] == "SET" {
			m[params[1]] = params[2]
			c.Write([]byte("OK\n"))
		} else if params[0] == "GET" {
			value := m[params[1]]
			if value == "" {
				c.Write([]byte("(nil)\n"))
			} else {
				c.Write([]byte(value + "\n"))
			}
		} else if params[0] == "CLOSE" {
			break
		} else if params[0] == "PING" {
			c.Write([]byte("PONG\n"))
		}
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

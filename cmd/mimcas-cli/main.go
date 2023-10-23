package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	var host string
	var port int
	flag.StringVar(&host, "host", "localhost", "Host to use for connection.")
	flag.IntVar(&port, "port", 20000, "Port to use for connection.")
	flag.Parse()


	c, err := net.Dial("tcp", host + ":" + strconv.Itoa(port))
	if err != nil {
			fmt.Println(err)
			return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		input, _ := reader.ReadString('\n')
		params := strings.Split(input, " ")
		response := ""
		switch params[0] { // each command needs a different processing of the response because it spects different lines
		case "set":
			// check > 3 params
			if len(params) < 3 {
				response = "Error: syntax for set is \"set <key> <value>\""
			} else {
				// prepare message to send to the server
				message := strings.Join(params, " ")
				fmt.Fprintf(c, message+"\n")
				response, _ = bufio.NewReader(c).ReadString('\n')
				// process response
			}
			//response = set(params)
		case "get":
			message := strings.Join(params, " ")
			fmt.Fprintf(c, message+"\n")
			response, _ = bufio.NewReader(c).ReadString('\n')
			// check = 2  params
			//response = get(params)
		case "mget":
			//check >= 2  params
			message := strings.Join(params, " ")
			fmt.Fprintf(c, message+"\n")
			response, _ = bufio.NewReader(c).ReadString('\n')
			//response = mget(params)
		case "del":
			message := strings.Join(params, " ")
			fmt.Fprintf(c, message+"\n")
			response, _ = bufio.NewReader(c).ReadString('\n')
			//response = delete(params)
		case "quit":
			message := strings.Join(params, " ")
			fmt.Fprintf(c, message+"\n")
			response, _ = bufio.NewReader(c).ReadString('\n')
			break
		case "ping":
			message := strings.Join(params, " ")
			fmt.Fprintf(c, message+"\n")
			response, _ = bufio.NewReader(c).ReadString('\n')
			//response = "pong\n"
		default:
			response = "ERR unknown command\n"
		}
		fmt.Print("->: " + response)
	}
}
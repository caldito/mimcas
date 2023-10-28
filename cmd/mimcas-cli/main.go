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
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("\nAvailable commands: set, get, mget, del, ping, quit")
	}
	flag.Parse()


	c, err := net.Dial("tcp", host + ":" + strconv.Itoa(port))
	if err != nil {
			fmt.Println(err)
			return
	}

	inputReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(c)

	for {
		fmt.Print(">> ")
		input, _ := inputReader.ReadString('\n')
		input = strings.TrimSpace(input)
		params := strings.Split(input, " ")
		response := ""
		switch params[0] { // each command needs a different processing of the response because it expects different lines
		case "set":
			if len(params) >= 3 {
				message := strings.Join(params, " ")
				fmt.Fprintf(c, message+"\n")
				response, _ = serverReader.ReadString('\n')
			} else {
				response = "Error: syntax for set is \"set <key> <value>\"\n"
			}
		case "get":
			if len(params) == 2 {
				message := strings.Join(params, " ")
				fmt.Fprintf(c, message+"\n")
				response, _ = serverReader.ReadString('\n')
				if response == "OK\n" {
					value, _:= serverReader.ReadString('\n')
					response += value
				}
			} else {
				response = "Error: syntax for get is \"get <key>\"\n"
			}
		case "mget":
			if len(params) >= 2 {
				message := strings.Join(params, " ")
				fmt.Fprintf(c, message+"\n")
				for _, _ = range params[1:] {
					status, _ := serverReader.ReadString('\n')
					response += status
					if status == "OK\n" {
						value, _:= serverReader.ReadString('\n')
						response += value
					}
				}		
			} else {
				response = "Error: syntax for mget is \"mget <key1> <key2> ...\"\n"
			}

		case "del":
			if len(params) == 2 {
				message := strings.Join(params, " ")
				fmt.Fprintf(c, message+"\n")
				response, _ = serverReader.ReadString('\n')
			} else {
				response = "Error: syntax for del is \"del <key>\"\n"
			}
		case "ping":
			if len(params) == 1 {
				message := strings.Join(params, " ")
				fmt.Fprintf(c, message+"\n")
				response, _ = serverReader.ReadString('\n')
			} else {
				response = "Error: syntax for ping is \"ping\"\n"
			}
		case "quit":
			break
		default:
			fmt.Println(params[0])
			response = "Error: unknown command\n"
		}
		fmt.Print(response)
	}
}

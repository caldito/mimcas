package main

import (
	// from other places
	"flag"
	"fmt"
	// "os"
	// "time"
)

func main() {
	//var repo string
	var port int
	//flag.StringVar(&repo, "repo", "", "url of the repository")
	flag.IntVar(&port, "port", 8080, "port to listen to")
	flag.Parse()
	// if repo == "" {
	// 	fmt.Println("Exiting, repo flag is not provided")
	// 	os.Exit(2)
	// }
	fmt.Println("Hello world")
}

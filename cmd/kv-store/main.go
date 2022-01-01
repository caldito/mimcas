package main

import (
	// from other places
	"net/http"
	"github.com/go-chi/chi/v5"
	"flag"
	"strconv"
	// "fmt"
	// "os"
	// "time"
)
var m = make(map[string]string)

func getValue(w http.ResponseWriter, r *http.Request) {
	keyParam := chi.URLParam(r, "key")
	w.Write([]byte(m[keyParam]))
}

func setValue(w http.ResponseWriter, r *http.Request) {
	keyParam := chi.URLParam(r, "key")
	valueParam := chi.URLParam(r, "val")
	m[keyParam] = valueParam
}

func main() {
	//var repo string
	var port int
	flag.IntVar(&port, "port", 8080, "port to listen to")
	flag.Parse()
	r := chi.NewRouter()
	r.Get("/value/{key}", getValue)
	r.Put("/value/{key}/{val}", setValue)
	
	http.ListenAndServe(":" + strconv.Itoa(port), r)
}

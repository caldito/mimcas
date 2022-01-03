package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"flag"
	"strconv"
	// "time"
	// "fmt"
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

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to listen to")
	flag.Parse()
	r := chi.NewRouter()
	r.Get("/keys/{key}", getKey)
	r.Put("/keys/{key}/{val}", setKey)
	r.Get("/health", getHealth)
	
	http.ListenAndServe(":" + strconv.Itoa(port), r)
}

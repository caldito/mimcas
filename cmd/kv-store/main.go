package main

import (
	// from other places
	"net/http"
	"github.com/go-chi/chi/v5"
	"flag"
	// 	"fmt"
	// "os"
	// "time"
)

func getArticle(w http.ResponseWriter, r *http.Request) {
	keyParam := chi.URLParam(r, "key")
	//slugParam := chi.URL(r, "slug")
	//article, err := database.GetArticle(date, slug)
  
	//if err != nil {
	//  w.WriteHeader(422)
	//  w.Write([]byte(fmt.Sprintf("error fetching article %s-%s: %v", dateParam, slugParam, err)))
	//  return
	//}
	//
	//if article == nil {
	//  w.WriteHeader(404)
	//  w.Write([]byte("article not found"))
	//  return
	//}
	w.Write([]byte("The key is: " + keyParam))
}

func main() {
	//var repo string
	var port string
	//flag.StringVar(&repo, "repo", "", "url of the repository")
	flag.StringVar(&port, "port", "8080", "port to listen to")
	flag.Parse()
	r := chi.NewRouter()
	r.Get("/value/{key}", getArticle)
	
	http.ListenAndServe(":" + port, r)
	// if repo == "" {
	// 	fmt.Println("Exiting, repo flag is not provided")
	// 	os.Exit(2)
	// }
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var name string

func main() {
	name = os.Getenv("POD_NAME")
	log.Println("protected app", name, "server ready")

	var connections int64
	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		atomic.AddInt64(&connections, 1)
		if connections > 1 {
			log.Println("Current connections greater than 1!", connections)
		}
	})

	http.HandleFunc("/disconnect", func(writer http.ResponseWriter, request *http.Request) {
		atomic.AddInt64(&connections, -1)
	})

	http.HandleFunc("/count", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "current connections: %d\n", connections)
	})
	http.ListenAndServe(":50001", nil)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(name, r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

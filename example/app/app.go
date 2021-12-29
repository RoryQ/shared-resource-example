package main

import (
	"log"
	"net/http"
	"os"
)

var name string

func main() {
	name = os.Getenv("POD_NAME")
	log.Println("takelock app", name, "server ready")

	http.Handle("/", loggingHandler(http.FileServer(http.Dir("/etc/podinfo"))))
	http.ListenAndServe(":50051", nil)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(name, r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

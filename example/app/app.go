package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var name string

func main() {
	name = os.Getenv("POD_NAME")
	log.Println("takelock app", name, "server ready")

	go useProtectedResource(context.Background())

	handler := loggingHandler(http.FileServer(http.Dir("/etc/podinfo")))
	http.Handle("/", handler)
	http.ListenAndServe(":50051", handler)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(name, r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

// this simulates using a protected resource i.e. only one process should access this resource at a time so a
// distributed lock is required to coordinate sharing access.
func useProtectedResource(ctx context.Context) {
	if err := connect(ctx); err != nil {
		panic(err)
	}

	log.Println(name, "connected")

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		c := http.Client{Timeout: 30 * time.Second}
		if _, err := c.Get("http://protected:50001/disconnect"); err != nil {
			log.Println(err)
		}
	}
}

func connect(ctx context.Context) error {
	c := http.Client{Timeout: 30 * time.Second}
	_, err := c.Get("http://protected:50001/connect")
	return err
}

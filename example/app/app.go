package main

import (
	"context"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/flowerinthenight/spindle"
)

var name string

func main() {
	name = os.Getenv("POD_NAME")
	log.Println("takelock app", name, "server ready")

	useProtectedResource(context.Background())

	//serveHTTP()
}

func useProtectedResource(ctx context.Context) {
	db, _ := spanner.NewClient(ctx, "projects/proj/instances/inst/databases/db")
	defer db.Close()

	if !tableExists(ctx, db, "ResourceLockTable") {
		panic("ResourceLockTable not found. Migrations have not been run")
	}

	lock := spindle.New(db, "ResourceLockTable", "protected-resource", withLeaseDuration(20*time.Second), spindle.WithId(name))
	lock.Run(ctx)

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			gotLock, token := lock.HasLock()
			if gotLock {
				log.Println("HasLock:", name, token)
				useSharedResource(ctx)

				// if lost shared resource lock then shutdown app
				return
			}
		}
	}
}

// this simulates using a protected resource i.e. only one process should access this resource at a time so a
// distributed lock is required to coordinate sharing access.
func useSharedResource(ctx context.Context) {
	c := http.Client{Timeout: 30 * time.Second}
	if _, err := c.Get("http://protected:50001/connect"); err != nil {
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

func serveHTTP() {
	loggingHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(name, r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
		})
	}
	handler := loggingHandler(http.FileServer(http.Dir("/etc/podinfo")))
	http.Handle("/", handler)
	http.ListenAndServe(":50051", handler)
}

func withLeaseDuration(d time.Duration) spindle.Option {
	return spindle.WithDuration(d.Milliseconds())
}

func tableExists(ctx context.Context, db *spanner.Client, tableName string) bool {
	existsStmt := spanner.NewStatement("SELECT table_name FROM information_schema.tables WHERE table_catalog = '' AND table_name = @table")
	existsStmt.Params["table"] = tableName
	ri := db.Single().Query(ctx, existsStmt)
	defer ri.Stop()
	_, err := ri.Next()
	return !(err == iterator.Done)
}

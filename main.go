package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flag"

	"github.com/gorilla/mux"
)

func main() {
	// We'll take three optional arguments -- the port, the dir to mount,
	// whether to populate a test dir.
	var baseDir = flag.String("baseDir", "./test", "The directory to mount the ls handler")
	var port = flag.Int("port", 8080, "The port to listen on")
	var test = flag.Bool("test", false, "Whether to artificially populate some files under ./test")
	flag.Parse()

	if *test {
		setupTestDir()
	}

	log.Printf("Starting router with base path %s on port %d", *baseDir, *port)

	lsHandler, err := NewLsHandler(*baseDir)
	if err != nil {
		log.Printf("Failed to create LS handler: %v", err)
		os.Exit(1)
	}

	// We only have one route, and that's /. Everything else we treat as the path.
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").Handler(lsHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", *port),
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// shunt the server to a goroutine -- mostly unecessary here, since we're
	// not doing anything else, but in case we wanted the server to do something
	// else as well.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Wait for SIGINT and SIGTERM for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	// Give the server 15 seconds to shut down
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server failed graceful shutdown:%+v", err)
	}

	if *test {
		cleanupTestDir()
	}

	log.Println("Shutting down")
	os.Exit(0)
}

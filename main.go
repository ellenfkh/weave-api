package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"flag"

	"github.com/gorilla/mux"
)

func main() {
	var baseDir = flag.String("baseDir", "./test", "The directory to mount the ls handler")
	var port = flag.Int("port", 8080, "The port to listen on")
	flag.Parse()

	log.Printf("Starting router with base path %s on port %d", *baseDir, *port)

	lsHandler, err := NewLsHandler(*baseDir)
	if err != nil {
		log.Printf("Failed ot create LS handler: %v", err)
		os.Exit(1)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").Handler(lsHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), router))
}

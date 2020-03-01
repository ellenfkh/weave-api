package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// http.Handler that just logs the path. For testing routes.
type LogHandler struct {
}

func (l *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, r.RequestURI)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "hello world"}`))
}

func main() {
	lsHandler := &LsHandler{baseDir: "./test"}
	//logHandler := LogHandler{}
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").Handler(lsHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

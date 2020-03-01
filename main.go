package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

type LsHandler struct {
	baseDir string
}

func (l *LsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := path.Join(l.baseDir, r.RequestURI)

	fi, err := os.Lstat(filename)
	if err != nil {
		writeError(w, err)
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsRegular():
		fmt.Println("regular file")

	case mode.IsDir():
		fmt.Println("directory")

	default:
		fmt.Println("unsupported filetype")
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

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
	lsHandler := &LsHandler{baseDir: "./"}
	//logHandler := LogHandler{}
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/").Handler(lsHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

// LsHandler serves HTTP requests by exploring the fs under baseDir/
type LsHandler struct {
	baseDir string
}

func NewLsHandler(baseDir string) (*LsHandler, error) {
	fi, err := os.Lstat(baseDir)
	if err != nil {
		return nil, errors.New("Error mounting basedir")
	}

	if !fi.Mode().IsDir() {
		return nil, errors.New("baseDir is not a directory")
	}

	return &LsHandler{baseDir: baseDir}, nil
}

// Write a 404 and log an error server-side
func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "{error: %+v}", err)
	log.Println(err)
}

// Send requests to cat or list, depending on whether it's a file or dir.
func (l *LsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename := path.Join(l.baseDir, r.RequestURI)

	fi, err := os.Lstat(filename)
	if err != nil {
		writeError(w, err)
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsRegular():
		data, err := l.cat(filename)
		if err != nil {
			writeError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		bytes, err := json.Marshal(&data)
		w.Write(bytes)
		return

	case mode.IsDir():
		files, err := l.list(filename)
		if err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		bytes, err := json.Marshal(files)
		if err != nil {
			writeError(w, err)
			return
		}
		w.Write(bytes)
		return

	default:
		writeError(w, errors.New("Unsupported filetype"))
		return
	}
}

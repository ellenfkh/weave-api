package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"syscall"
)

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Println(err)
}

type fileInfo struct {
	name        string
	permissions string
	owner       string
	size        int64
}

func (l *LsHandler) list(filename string) ([]*fileInfo, error) {
	files, err := ioutil.ReadDir(filename)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var fileArray []*fileInfo

	for _, file := range files {
		// fIXME: this assumes posix
		fStat := file.Sys().(*syscall.Stat_t)
		owner, _ := user.LookupId(fmt.Sprint(fStat.Uid))

		fileArray = append(fileArray, &fileInfo{
			name:        file.Name(),
			permissions: file.Mode().String(),
			owner:       owner.Name,
			size:        file.Size(),
		})
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return fileArray, nil
}

func (l *LsHandler) cat(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Print(string(data))

	return data, nil
}

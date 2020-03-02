package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"syscall"
)

// List the files in the given dir, and return a slice of fileInfos
// describing the contents of the dir. Returns an error if we can't read the
// dir, but only logs and omits if we can't stat a file
func (l *LsHandler) list(filename string) ([]*fileInfo, error) {
	files, err := ioutil.ReadDir(filename)
	if err != nil {
		return nil, err
	}

	var fileArray []*fileInfo

	for _, file := range files {
		// don't fail if one file throws, just log
		stat, err := extractInfo(file)
		if err != nil {
			log.Printf("Failed to stat file %s: %v ", file.Name(), err)
		} else {
			fileArray = append(fileArray, stat)
		}
	}

	return fileArray, nil
}

// Cat the contents of a given file, return its contents along with its stat.
func (l *LsHandler) cat(filename string) (*fileContents, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	file, err := os.Lstat(filename)
	if err != nil {
		return nil, err
	}

	stat, err := extractInfo(file)

	if err != nil {
		return nil, err
	}

	return &fileContents{Stat: stat, Contents: string(data)}, nil
}

// utility to massage an os.FileInfo into our custom fileInfo struct, stripping
// the pieces we don't care about
func extractInfo(file os.FileInfo) (*fileInfo, error) {
	// FIXME: this has been tested only on Mac, and runs into the problem where
	// in Docker for Mac, the owner of the mount is the user that the container
	// runs as (which is usually root or docker), depending on how consciencious
	// the host user is. Sys is not necessarily os-agnostic, so this may be
	// broken on Linux hosts.
	fStat := file.Sys().(*syscall.Stat_t)
	owner, err := user.LookupId(fmt.Sprint(fStat.Uid))

	if err != nil {
		return nil, err
	}

	return &fileInfo{
		Name:        file.Name(),
		Permissions: file.Mode().String(),
		Owner:       owner.Name,
		Size:        file.Size(),
		IsDir:       file.IsDir(),
	}, nil
}

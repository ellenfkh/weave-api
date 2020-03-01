package main

type fileContents struct {
	Stat     *fileInfo
	Contents string
}

type fileInfo struct {
	Name        string
	Permissions string
	Owner       string
	Size        int64
	IsDir       bool
}

type files []fileInfo

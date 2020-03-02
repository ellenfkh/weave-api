package main

// These structs describe the objects returned by the API. The fs information is
// marshaled into files or fileInfo blobs in the response.

type fileInfo struct {
	Name        string
	Permissions string
	Owner       string
	Size        int64
	IsDir       bool
}

type fileContents struct {
	Stat     *fileInfo
	Contents string
}

type files []fileInfo

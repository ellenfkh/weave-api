package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func setupTestDir() {
	// Set up a test dir that looks like this:
	// test
	//		|- file
	//		|- .hidden
	// 		|- dir/
	// 			|- nested
	// 		|- hidden_dir/
	// 			|- nested
	// 		|- empty/

	os.Mkdir("./test", os.FileMode(0755))
	ioutil.WriteFile("./test/file", []byte("test-file"), 0644)
	ioutil.WriteFile("./test/.hidden", []byte("test-hidden"), 0644)

	os.Mkdir("./test/dir", os.FileMode(0755))
	os.Mkdir("./test/.hidden_dir", os.FileMode(0755))
	os.Mkdir("./test/empty", os.FileMode(0755))

	ioutil.WriteFile("./test/dir/nested", []byte("test-dir-nested"), 0644)
	ioutil.WriteFile("./test/.hidden_dir/nested", []byte("test-hidden_dir-nested"), 0644)

	fmt.Println("Created test dir")
}

func cleanupTestDir() {
	os.RemoveAll("./test")
	fmt.Println("Removed test dir")
}

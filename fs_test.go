package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Wrap tests in helper funcs to create a nice test dir, then nukeit
func TestMain(m *testing.M) {
	setupTestDir()

	m.Run()
	// FIXME: don't clean up for now, since this is useful for testing the handler
	// cleanupTestDir()
}

func TestCat(t *testing.T) {
	lsHandler := LsHandler{}
	user, err := user.Current()
	if err != nil {
		t.Errorf("Failed to get current user")
	}

	// Should fail if we try to cat a dir
	data, err := lsHandler.cat("./test")
	assert.Error(t, err, "Should have thrown an error, this is a dir")
	assert.Nil(t, data, "Should have failed to cat a dir")

	// Should fail if we try to cat a dir
	data, err = lsHandler.cat("./fake")
	assert.Error(t, err, "Should have thrown an error, nonexiestent")
	assert.Nil(t, data, "Should have cat a nonexistent file")

	// Get a normal file's contents

	data, err = lsHandler.cat("./test/file")
	assert.Equal(t, "test-file", data.Contents, "Should have thrown an error, nonexiestent")
	assert.Equal(t,
		&fileInfo{
			Name:        "file",
			Permissions: os.FileMode(0644).String(),
			Owner:       user.Name,
			Size:        int64(len([]byte("test-file"))),
			IsDir:       false,
		},
		data.Stat, "Failed to cat file")
	assert.Nil(t, err, "Should not have thrown an error")
}

func TestList(t *testing.T) {
	lsHandler := LsHandler{}

	// Should get the right number of files, including dirs and hiddens
	files, err := lsHandler.list("./test")
	assert.Equal(t, 4, len(files),
		"Found wrong number of files in ./test")
	assert.Nil(t, err, "Unexpected error listing ./test")

	// Should return an error if we try to list a bad path
	files, err = lsHandler.list("./fake")
	assert.Error(t, err, "Should have thrown an error nonexistent dir")
	assert.Nil(t, files, "Should have failed to list nonexisting dir")

	// Should return an error if we try to list a file
	files, err = lsHandler.list("./test/file")
	assert.Error(t, err, "Should have thrown an error nonexistent dir")
	assert.Nil(t, files, "Should have failed to list nonexisting dir")

	// Should be able to walk compound paths
	files, err = lsHandler.list("./test/dir")
	assert.Equal(t, 1, len(files),
		"Found wrong number of files in ./test/dir")
	assert.Nil(t, err, "Unexpected error listing ./test/dir")

	// Check that the file stats are as expected
	user, err := user.Current()

	if err != nil {
		t.Errorf("Failed to get current user")
	}

	assert.Equal(t,
		&fileInfo{
			Name:        "nested",
			Permissions: os.FileMode(0644).String(),
			Owner:       user.Name,
			Size:        int64(len([]byte("test-dir-nested"))),
			IsDir:       false,
		},
		files[0],
		"Got malformed file stats")
}

func setupTestDir() {
	// Set up a test dir that looks like this:
	// test
	//		|- file
	//		|- .hidden
	// 		|- dir/
	// 			|- nested
	// 		|- hidden_dir/
	// 			|- nested

	os.Mkdir("./test", os.FileMode(0755))
	ioutil.WriteFile("./test/file", []byte("test-file"), 0644)
	ioutil.WriteFile("./test/.hidden", []byte("test-hidden"), 0644)

	// os.Create("./test/.hidden")

	os.Mkdir("./test/dir", os.FileMode(0755))
	os.Mkdir("./test/.hidden_dir", os.FileMode(0755))

	ioutil.WriteFile("./test/dir/nested", []byte("test-dir-nested"), 0644)
	ioutil.WriteFile("./test/.hidden_dir/nested", []byte("test-hidden_dir-nested"), 0644)

	fmt.Println("Created test dir")
}

func cleanupTestDir() {
	os.RemoveAll("./test")
	fmt.Println("Removed test dir")
}

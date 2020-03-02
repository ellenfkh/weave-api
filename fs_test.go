package main

import (
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Wrap tests in helper funcs to create a nice test dir, then nukeit
func TestMain(m *testing.M) {
	setupTestDir()
	m.Run()
	cleanupTestDir()
}

func TestCat(t *testing.T) {
	// We don't care about the baseDir of lsHandler -- we're only testing the fs
	// operations, not the path handling
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
	// We don't care about the baseDir of lsHandler -- we're only testing the fs
	// operations, not the path handling
	lsHandler := LsHandler{}
	user, err := user.Current()
	if err != nil {
		t.Errorf("Failed to get current user")
	}

	// Should get the right number of files, including dirs and hiddens
	files, err := lsHandler.list("./test")
	assert.Equal(t, 5, len(files),
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

	// Handle empty dirs
	files, err = lsHandler.list("./test/empty")
	assert.Equal(t, 0, len(files),
		"Found wrong number of files in ./test/empty")
	assert.Nil(t, err, "Unexpected error listing ./test/empty")

	// Should be able to walk compound paths
	files, err = lsHandler.list("./test/dir")
	assert.Equal(t, 1, len(files),
		"Found wrong number of files in ./test/dir")
	assert.Nil(t, err, "Unexpected error listing ./test/dir")

	// Check that the file stats are as expected
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

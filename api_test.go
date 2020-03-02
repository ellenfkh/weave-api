package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testNewLsHAndler(t *testing.T) {
	handler, err := NewLsHandler("Fake-Direcotry")
	assert.Error(t, err, "Should have thrown an error, is not a dir")
	assert.Nil(t, handler, "Should have failed to create a handler")
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("Something bad happened!")

	writeError(w, err)

	assert.Equal(t, fmt.Sprintf("{error: %+v}", err), w.Body.String())
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}

// This is just sanity checking, not checking for correctness -- most of the
// logical testing is done in fs_test. Hwoever, there is a case here where we're
// not explicitly making sure that files are catted and dirs are ls-ed. We're
// relying on the assumption fhat if that weren't happening, we'd be getting an
// error from the os package for trying to do illegal operations.
func TestServeHttp(t *testing.T) {
	setupTestDir()
	handler, _ := NewLsHandler("./test")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/fake", nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/file", nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/dir", nil)

	handler.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
}

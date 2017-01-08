package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	QueryHandler().ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusOK)
	assert.Contains(w.Body.String(), `<input type="text" name="q">`)
}

func TestOpenSearch(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/opensearch.xml", nil)
	w := httptest.NewRecorder()

	OpenSearchHandler().ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusOK)
	assert.Contains(w.Body.String(), "<OpenSearchDescription")
}

func TestCommand(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/?q=ping", nil)
	w := httptest.NewRecorder()

	QueryHandler().ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusOK)

	body := w.Body.String()
	tokens := strings.Split(body, " ")
	assert.Len(tokens, 2)
	assert.Equal(tokens[0], "pong")

	n, err := strconv.ParseInt(tokens[1], 10, 64)
	assert.Nil(err)
	assert.WithinDuration(time.Now(), time.Unix(n, 0), 1*time.Second)
}

func TestInvalidCommand(t *testing.T) {
	assert := assert.New(t)

	db, _ = bolt.Open("test.db", 0600, nil)
	defer db.Close()

	r, _ := http.NewRequest("GET", "/?q=asdf", nil)
	w := httptest.NewRecorder()

	QueryHandler().ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusBadRequest)

	body := w.Body.String()
	assert.Equal(body, "Invalid Command: asdf\n")
}

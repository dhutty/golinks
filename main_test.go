package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

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

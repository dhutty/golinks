package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

type Explode struct{}

func (e Explode) Name() string {
	return "explode"
}

func (e Explode) Desc() string {
	return `explode

	This command just quite literally explodes and always returns an error!
	`
}

func (e Explode) Exec(w http.ResponseWriter, args []string) error {
	return errors.New("kaboom")
}

func TestRender(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	templates.Load()

	render(w, "index", nil)

	assert.Equal(w.Code, http.StatusOK)
	assert.Contains(w.Body.String(), `<input type="text" name="q">`)
}

func TestRenderError(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()

	templates.Load()

	render(w, "asdf", nil)

	assert.Equal(w.Code, http.StatusInternalServerError)
}

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	templates.Load()

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

func TestCommandError(t *testing.T) {
	assert := assert.New(t)

	db, _ = bolt.Open("test.db", 0600, nil)
	defer db.Close()

	RegisterCommand("explode", Explode{})

	r, _ := http.NewRequest("GET", "/?q=explode", nil)
	w := httptest.NewRecorder()

	QueryHandler().ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusInternalServerError)

	body := w.Body.String()
	assert.Equal(body, "Error processing command explode: kaboom\n")
}

func TestCommandBookmark(t *testing.T) {
	assert := assert.New(t)

	db, _ = bolt.Open("test.db", 0600, nil)
	defer db.Close()

	EnsureDefaultBookmarks()

	r, _ := http.NewRequest("GET", "/?q=g%20foo%20bar", nil)
	w := httptest.NewRecorder()

	QueryHandler().ServeHTTP(w, r)
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})

	assert.Equal(
		w.Header().Get("Location"),
		"https://www.google.com/search?q=foo bar&btnK",
	)
}

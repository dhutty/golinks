package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroBookmark(t *testing.T) {
	assert := assert.New(t)

	bookmark := Bookmark{}
	assert.Equal(bookmark.Name(), "")
	assert.Equal(bookmark.URL(), "")

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	bookmark.Exec(w, r, "")
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})
}

func TestBookmark(t *testing.T) {
	assert := assert.New(t)

	bookmark := Bookmark{
		name: "g",
		url:  "https://www.google.com/search?q=%s&btnK",
	}
	assert.Equal(bookmark.Name(), "g")
	assert.Equal(bookmark.URL(), "https://www.google.com/search?q=%s&btnK")

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	q := "foo bar"
	bookmark.Exec(w, r, q)
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})

	assert.Equal(
		w.Header().Get("Location"),
		fmt.Sprintf(
			"https://www.google.com/search?q=%s&btnK",
			q,
		),
	)
}

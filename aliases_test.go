package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroAlias(t *testing.T) {
	assert := assert.New(t)

	alias := Alias{}
	assert.Equal(alias.Name(), "")
	assert.Equal(alias.URL(), "")

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	err := alias.Exec(w, r, "")
	assert.Nil(err)
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})
}

func TestAlias(t *testing.T) {
	assert := assert.New(t)

	alias := Alias{name: "g", url: "https://www.google.com/search?q=%s&btnK"}
	assert.Equal(alias.Name(), "g")
	assert.Equal(alias.URL(), "https://www.google.com/search?q=%s&btnK")

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	q := "foo bar"
	err := alias.Exec(w, r, q)
	assert.Nil(err)
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

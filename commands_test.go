package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

// Foo ...
type Foo struct{}

// Name ...
func (f Foo) Name() string {
	return "foo"
}

// Desc ...
func (f Foo) Desc() string {
	return "foo bar"
}

// Exec ...
func (f Foo) Exec(w http.ResponseWriter, args []string) error {
	w.Write([]byte(fmt.Sprintf("foo bar")))
	return nil
}

func TestNewCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := Foo{}
	assert.Equal(cmd.Name(), "foo")
	assert.Equal(cmd.Desc(), "foo bar")

	w := httptest.NewRecorder()
	args := []string{}
	err := cmd.Exec(w, args)
	assert.Nil(err)

	body := w.Body.String()
	assert.Equal(body, "foo bar")
}

func TestPingCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := Ping{}
	assert.Equal(cmd.Name(), "ping")
	assert.Contains(cmd.Desc(), "ping")

	w := httptest.NewRecorder()
	args := []string{}
	err := cmd.Exec(w, args)
	assert.Nil(err)

	body := w.Body.String()
	assert.Regexp("^pong [0-9]+$", body)
}

func TestListCommand(t *testing.T) {
	assert := assert.New(t)

	db, _ = bolt.Open("test.db", 0600, nil)
	defer db.Close()

	EnsureDefaultBookmarks()

	cmd := List{}
	assert.Equal(cmd.Name(), "list")
	assert.Contains(cmd.Desc(), "list")

	w := httptest.NewRecorder()
	args := []string{}
	err := cmd.Exec(w, args)
	assert.Nil(err)

	body := w.Body.String()
	cmds := strings.Split(body, "\n")
	assert.Contains(cmds, "ping")
	assert.Contains(cmds, "list")
	assert.Contains(cmds, "help")
	assert.Contains(cmds, "g")
	assert.Contains(cmds, "yt")
	assert.Contains(cmds, "go")
}

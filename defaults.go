package main

import (
	"net/http/httptest"
)

// DefaultBookmarks ...
var DefaultBookmarks map[string]string

func init() {
	DefaultBookmarks = map[string]string{
		"g":  "https://www.google.com/search?q=%s&btnK",
		"gl": "https://www.google.com/search?q=%s&btnI",
		"gh": "https://github.com/search?q=%s&ref=opensearch",
		"go": "https://golang.org/search?q=%s",
	}
}

// EnsureDefaultBookmarks ...
func EnsureDefaultBookmarks() error {
	for k, v := range DefaultBookmarks {
		if _, ok := LookupBookmark(k); !ok {
			w := httptest.NewRecorder()
			args := []string{k, v}
			add := Add{}
			if err := add.Exec(w, args); err != nil {
				return err
			}
		}
	}
	return nil
}

package main

import (
	"net/http/httptest"
)

// DefaultAliases ...
var DefaultAliases map[string]string

func init() {
	DefaultAliases = map[string]string{
		"g":  "https://www.google.com/search?q=%s&btnI",
		"gh": "https://github.com/search?q=%s&ref=opensearch",
	}
}

// EnsureDefaultAliases ...
func EnsureDefaultAliases() error {
	for k, v := range DefaultAliases {
		if _, ok := LookupAlias(k); !ok {
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

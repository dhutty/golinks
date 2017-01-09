package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

// Bookmark ...
type Bookmark struct {
	name string
	url  string
}

// Name ...
func (b Bookmark) Name() string {
	return b.name
}

// URL ...
func (b Bookmark) URL() string {
	return b.url
}

// Exec ...
func (b Bookmark) Exec(w http.ResponseWriter, r *http.Request, q string) {
	url := fmt.Sprintf(b.url, q)
	http.Redirect(w, r, url, http.StatusFound)
}

// LookupBookmark ...
func LookupBookmark(name string) (bookmark Bookmark, ok bool) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bookmarks"))
		if b == nil {
			return nil
		}

		v := b.Get([]byte(name))
		if v != nil {
			bookmark.name = name
			bookmark.url = string(v)
			ok = true
		}

		return nil
	})

	if err != nil {
		log.Printf("error looking up bookmark for %s: %s", name, err)
	}

	return
}

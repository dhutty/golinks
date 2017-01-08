package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

type Alias struct {
	name string
	url  string
}

// Name ...
func (a Alias) Name() string {
	return a.name
}

// Exec ...
func (a Alias) Exec(w http.ResponseWriter, r *http.Request, q string) error {
	url := fmt.Sprintf(a.url, q)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func LookupAlias(name string) (alias Alias, ok bool) {
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("aliases"))
		if b == nil {
			return nil
		}

		v := b.Get([]byte(name))
		if v != nil {
			alias.name = name
			alias.url = string(v)
			ok = true
		}

		return nil
	})

	if err != nil {
		log.Printf("error looking up alias for %s: %s", name, err)
	}

	return
}

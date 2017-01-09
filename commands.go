package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

// Command ...
type Command interface {
	Name() string
	Desc() string
	Exec(w http.ResponseWriter, args []string) error
}

var commands map[string]Command

func init() {
	commands = make(map[string]Command)
	RegisterCommand("list", List{})
	RegisterCommand("ping", Ping{})
	RegisterCommand("help", Help{})
	RegisterCommand("date", Date{})
	RegisterCommand("time", Time{})
	RegisterCommand("add", Add{})
}

// RegisterCommand ...
func RegisterCommand(name string, command Command) {
	commands[name] = command
}

// LookupCommand ...
func LookupCommand(name string) Command {
	command, ok := commands[name]
	if ok {
		return command
	}
	return nil
}

// Ping ...
type Ping struct{}

// Name ...
func (p Ping) Name() string {
	return "ping"
}

// Desc ...
func (p Ping) Desc() string {
	return `ping

	Responds with "pong <ts>" where ts is the current UNIX timestamp.
	`
}

// Exec ...
func (p Ping) Exec(w http.ResponseWriter, args []string) error {
	w.Write([]byte(fmt.Sprintf("pong %d", time.Now().Unix())))
	return nil
}

// List ...
type List struct{}

// Name ...
func (p List) Name() string {
	return "list"
}

// Desc ...
func (p List) Desc() string {
	return `list

	Lists all available commands.
	`
}

// Exec ...
func (p List) Exec(w http.ResponseWriter, args []string) error {
	var bs, cs []string

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("bookmarks"))
		if b == nil {
			return nil
		}

		err := b.ForEach(func(k, v []byte) error {
			bs = append(bs, string(k))
			return nil
		})

		return err
	})
	if err != nil {
		log.Printf("error reading list of bookmarks: %s", err)
	}

	for k := range commands {
		cs = append(cs, k)
	}

	w.Write(
		[]byte(
			fmt.Sprintf(
				"Bookmarks:\n\n%s\n\nCommands:\n\n%s",
				strings.Join(bs, "\n"),
				strings.Join(cs, "\n"),
			),
		),
	)

	return nil
}

// Help ...
type Help struct{}

// Name ...
func (p Help) Name() string {
	return "help"
}

// Desc ...
func (p Help) Desc() string {
	return `help

	Display general helpful information
	`
}

// Exec ...
func (p Help) Exec(w http.ResponseWriter, args []string) error {
	render(w, "help", nil)
	return nil
}

// Date ...
type Date struct{}

// Name ...
func (p Date) Name() string {
	return "date"
}

// Desc ...
func (p Date) Desc() string {
	return `date

	Display the current date and time
	`
}

// Exec ...
func (p Date) Exec(w http.ResponseWriter, args []string) error {
	w.Write([]byte(time.Now().Format(http.TimeFormat)))
	return nil
}

// Time ...
type Time struct{}

// Name ...
func (p Time) Name() string {
	return "time"
}

// Desc ...
func (p Time) Desc() string {
	return `time

	Display the current time
	`
}

// Exec ...
func (p Time) Exec(w http.ResponseWriter, args []string) error {
	w.Write([]byte(time.Now().Format("15:04:05")))
	return nil
}

// Add ...
type Add struct{}

// Name ...
func (p Add) Name() string {
	return "add"
}

// Desc ...
func (p Add) Desc() string {
	return `add [name] [url]

	Adds a new bookmark with the given name that will redierct to the given
	url passing arguments as %s. For example:

	add g http://google.com/search?btnK&q=%s

	Will add a new command called 'g' which will redirect to Google's search
	passing in arguments as '%s'
	`
}

// Exec ...
func (p Add) Exec(w http.ResponseWriter, args []string) error {
	var name, url string

	if len(args) == 2 {
		name, url = args[0], args[1]
	} else {
		return fmt.Errorf("expected 2 arguments got %d", len(args))
	}

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("bookmarks"))
		if err != nil {
			log.Printf("create bucket failed: %s", err)
			return err
		}

		err = b.Put([]byte(name), []byte(url))
		if err != nil {
			log.Printf("put key failed: %s", err)
			return err
		}

		w.Write([]byte("OK"))

		return nil
	})

	return err
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{}
	assert.Equal(cfg.Title, "")
	assert.Equal(cfg.FQDN, "")
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	cfg := Config{Title: "foo", FQDN: "bar.com"}
	assert.Equal(cfg.Title, "foo")
	assert.Equal(cfg.FQDN, "bar.com")
}

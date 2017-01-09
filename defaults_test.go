package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureDefaultBookmarksError(t *testing.T) {
	err := EnsureDefaultBookmarks()
	assert.NotNil(t, err)
}

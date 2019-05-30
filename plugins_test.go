package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err)
	os.Setenv("ASDF_DATA_DIR", path.Join(cwd, "fixtures", "asdf"))

	plugins, err := List()
	assert.Equal(t, []string{"go", "python", "ruby"}, plugins)

	os.Unsetenv("ASDF_DATA_DIR")
}

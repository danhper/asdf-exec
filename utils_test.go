package main

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLines(t *testing.T) {
	assert.Equal(t, []string{"foo", "bar"}, ReadLines("foo\nbar"))
	assert.Equal(t, []string{"foo", "bar"}, ReadLines("# hello\nfoo #bar \n  bar  "))
}

func TestRemoveAsdfPath(t *testing.T) {
	asdfShimPath := path.Join(GetAsdfDataPath(), "shims")
	home := os.Getenv("HOME")
	homeBin := path.Join(home, "bin")

	currentPath := []string{homeBin, asdfShimPath, "/usr/bin", "/bin"}
	actualPath := RemoveAsdfPath(strings.Join(currentPath, ":"))
	expectedPath := strings.Join([]string{homeBin, "/usr/bin", "/bin"}, ":")
	assert.Equal(t, expectedPath, actualPath)
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadLines(t *testing.T) {
	assert.Equal(t, []string{"foo", "bar"}, ReadLines("foo\nbar"))
	assert.Equal(t, []string{"foo", "bar"}, ReadLines("# hello\nfoo #bar \n  bar  "))
}

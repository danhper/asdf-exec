package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBool(t *testing.T) {
	for _, value := range []string{"yes", "1", "true", "True", "YES"} {
		res, err := ParseBool(value)
		assert.Nil(t, err)
		assert.Equal(t, true, res)
	}
	for _, value := range []string{"no", "0", "false", "False", "NO"} {
		res, err := ParseBool(value)
		assert.Nil(t, err)
		assert.Equal(t, false, res)
	}
	for _, value := range []string{"aaq", "something crazy", "123"} {
		_, err := ParseBool(value)
		assert.NotNil(t, err)
	}
}

func TestConfigFromString(t *testing.T) {
	config, err := ConfigFromString(`
	# single line comment
	legacy_version_file = yes # end of line comment
	# empty line

	use_release_candidates = true
	unknown_key = false
	`)
	assert.Nil(t, err)
	assert.True(t, config.LegacyVersionFile)

	_, err = ConfigFromString(`
	# bad format
	this is a bad line
	`)
	assert.NotNil(t, err)
}

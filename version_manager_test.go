package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	assert.Equal(t, []string{"3.7.2"}, ParseVersion(" 3.7.2"))
	assert.Equal(t, []string{"3.7.2", "2.7.16", "system"}, ParseVersion(" 3.7.2 2.7.16 system"))
}

func TestFindVersionsInEnv(t *testing.T) {
	_, found := FindVersionsInEnv("python")
	assert.False(t, found)
	os.Setenv("ASDF_PYTHON_VERSION", "3.7.2 2.7.16")
	versions, found := FindVersionsInEnv("python")
	assert.True(t, found)
	assert.Len(t, versions, 2)
	assert.Equal(t, []string{"3.7.2", "2.7.16"}, versions)
	os.Unsetenv("ASDF_PYTHON_VERSION")
}

func TestFindVersionsInToolFileContent(t *testing.T) {
	content := `
	# some comments
	python 3.6.7 2.7.11 system  # fallback to system

	ruby 2.6.2
	`
	versions, found := FindVersionsInToolFileContent("python", content)
	assert.True(t, found)
	assert.Len(t, versions, 3)
	assert.Equal(t, []string{"3.6.7", "2.7.11", "system"}, versions)

	versions, found = FindVersionsInToolFileContent("ruby", content)
	assert.True(t, found)
	assert.Len(t, versions, 1)
	assert.Equal(t, []string{"2.6.2"}, versions)

	_, found = FindVersionsInToolFileContent("nodejs", content)
	assert.False(t, found)
}

func TestFindVersions(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err)
	config := Config{LegacyVersionFile: false}

	os.Setenv("HOME", "/tmp")
	assert.Nil(t, os.Chdir("/tmp"))
	_, found, err := FindVersions("python", config)
	assert.Nil(t, err)
	assert.False(t, found)

	assert.Nil(t, os.Chdir(path.Join(cwd, "fixtures", "some-dir", "nested-dir")))
	versions, found, err := FindVersions("python", config)
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, []string{"3.6.7", "2.7.11"}, versions)

	os.Setenv("HOME", path.Join(cwd, "fixtures", "some-dir"))
	assert.Nil(t, os.Chdir("/tmp"))
	versions, found, err = FindVersions("python", config)
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, []string{"3.6.7", "2.7.11"}, versions)

	assert.Nil(t, os.Chdir(cwd))
}

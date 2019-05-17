package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseExecutableLine(t *testing.T) {
	executable, err := ParseExecutableLine("python", "# asdf-plugin: python 3.6.7")
	assert.Nil(t, err)
	assert.Equal(t, executable.PluginName, "python")
	assert.Equal(t, executable.PluginVersion, "3.6.7")
}

func TestGetExecutablesFromShim(t *testing.T) {
	shimContent := `
	#!/usr/bin/env bash
    # asdf-plugin: python 3.6.7
	# asdf-plugin: python 2.7.11
	exec /home/daniel/.asdf/bin/asdf exec "python" "$@"
	`
	executables, err := GetExecutablesFromShim("python", shimContent)
	assert.Nil(t, err)
	assert.Len(t, executables, 2)
}

func TestFindExecutable(t *testing.T) {
	config := Config{LegacyVersionFile: false}
	cwd, err := os.Getwd()
	assert.Nil(t, err)
	shimPath := path.Join(cwd, "fixtures", "asdf", "shims", "flask")

	os.Setenv("HOME", "/tmp")
	assert.Nil(t, os.Chdir("/tmp"))

	_, found, err := FindExecutable(shimPath, config)
	assert.Nil(t, err)
	assert.False(t, found)

	assert.Nil(t, os.Chdir(path.Join(cwd, "fixtures", "some-dir", "nested-dir")))
	executable, found, err := FindExecutable(shimPath, config)
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, "python", executable.PluginName)
	assert.Equal(t, "3.6.7", executable.PluginVersion)

	assert.Nil(t, os.Chdir(cwd))
}

func TestGetExecutablePath(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err)
	os.Setenv("ASDF_DATA_DIR", path.Join(cwd, "fixtures", "asdf"))

	executable := Executable{Name: "2to3", PluginName: "python", PluginVersion: "2.7.11"}
	executablePath, err := GetExecutablePath(executable)
	assert.Nil(t, err)
	expected := path.Join(os.Getenv("ASDF_DATA_DIR"), "installs", "python", "2.7.11", "bin", "2to3")
	assert.Equal(t, expected, executablePath)

	// check it works with list-bin-paths
	executable = Executable{Name: "go", PluginName: "go", PluginVersion: "1.9.1"}
	executablePath, err = GetExecutablePath(executable)
	assert.Nil(t, err)
	expected = path.Join(os.Getenv("ASDF_DATA_DIR"), "installs", "go", "1.9.1", "go", "bin", "go")
	assert.Equal(t, expected, executablePath)

	os.Unsetenv("ASDF_DATA_DIR")
}

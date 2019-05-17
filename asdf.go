package main

import (
	"os"
	"path"
)

// GetAsdfDataPath returns the path for asdf
func GetAsdfDataPath() string {
	dir := os.Getenv("ASDF_DATA_DIR")
	if dir != "" {
		return dir
	}
	return path.Join(os.Getenv("HOME"), ".asdf")
}

// GetPluginPath returns the path of the plugin
func GetPluginPath(plugin string) string {
	return path.Join(GetAsdfDataPath(), "plugins", plugin)
}

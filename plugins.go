package main

import (
	"io/ioutil"
	"path"
)

func getPluginDir() string {
	dir := GetAsdfDataPath()
	return path.Join(dir, "plugins")
}

// List returns the list of installed plugins
func List() (plugins []string, err error) {
	dirs, err := ioutil.ReadDir(getPluginDir())
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			plugins = append(plugins, dir.Name())
		}
	}
	return
}

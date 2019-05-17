package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	asdfPluginPrefix  string = "# asdf-plugin: "
	listBinPathScript string = "list-bin-paths"
)

// Executable is an instance of a single executable
type Executable struct {
	Name          string
	PluginName    string
	PluginVersion string
}

// ParseExecutableLine returns an executable from a shim plugin line
func ParseExecutableLine(name string, fullLine string) (Executable, error) {
	line := strings.Replace(fullLine, asdfPluginPrefix, "", -1)
	tokens := strings.Split(line, " ")
	if len(tokens) != 2 {
		return Executable{}, fmt.Errorf("bad line %s", fullLine)
	}
	return Executable{
		Name:          name,
		PluginName:    strings.TrimSpace(tokens[0]),
		PluginVersion: strings.TrimSpace(tokens[1]),
	}, nil
}

// GetExecutablesFromShim retrieves all the executable for a shim
func GetExecutablesFromShim(name string, content string) (executables []Executable, err error) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, asdfPluginPrefix) {
			executable, err := ParseExecutableLine(name, line)
			if err != nil {
				return executables, err
			}
			executables = append(executables, executable)
		}
	}
	return
}

// GetExecutablesFromShimFile retrieves all the executable for a shim file
func GetExecutablesFromShimFile(filepath string) ([]Executable, error) {
	name := path.Base(filepath)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []Executable{}, err
	}
	return GetExecutablesFromShim(name, string(content))
}

// FindExecutable returns the path to the executable to be executed
func FindExecutable(filepath string, config Config) (Executable, bool, error) {
	executables, err := GetExecutablesFromShimFile(filepath)
	if err != nil {
		return Executable{}, false, err
	}

	plugins := make(map[string][]Executable)

	for _, executable := range executables {
		pluginExecutables, ok := plugins[executable.PluginName]
		if !ok {
			pluginExecutables = []Executable{}
		}
		pluginExecutables = append(pluginExecutables, executable)
		plugins[executable.PluginName] = pluginExecutables
	}

	for plugin, pluginExecutables := range plugins {
		toolVersions, found, err := FindVersions(plugin, config)
		if err != nil {
			return Executable{}, false, err
		}
		if !found {
			continue
		}
		for _, toolVersion := range toolVersions {
			for _, executable := range pluginExecutables {
				if toolVersion == executable.PluginVersion {
					return executable, true, nil
				}
			}
		}
	}

	return Executable{}, false, nil
}

// GetShimPath returns the path of the shim
func GetShimPath(shimName string) string {
	return path.Join(GetAsdfDataPath(), "shims", shimName)
}

// GetExecutablePath returns the path of the executable
func GetExecutablePath(executable Executable) (string, error) {
	pluginPath := GetPluginPath(executable.PluginName)
	installPath := path.Join(
		GetAsdfDataPath(),
		"installs",
		executable.PluginName,
		executable.PluginVersion,
	)

	listBinPath := path.Join(pluginPath, "bin", listBinPathScript)
	if _, err := os.Stat(listBinPath); err != nil {
		exePath := path.Join(installPath, "bin", executable.Name)
		if _, err := os.Stat(exePath); err != nil {
			return "", fmt.Errorf("executable not found")
		}
		return exePath, nil
	}

	rawBinPaths, err := exec.Command("bash", listBinPath).Output()
	if err != nil {
		return "", err
	}
	paths := strings.Split(string(rawBinPaths), " ")
	for _, binPath := range paths {
		binPath = strings.TrimSpace(binPath)
		exePath := path.Join(installPath, binPath, executable.Name)
		if _, err := os.Stat(exePath); err == nil {
			return exePath, nil
		}
	}
	return "", fmt.Errorf("executable not found")
}

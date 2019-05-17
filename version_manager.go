package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const toolVersionsFile string = ".tool-versions"

// ParseVersion parses the raw version
func ParseVersion(rawVersions string) []string {
	var versions []string
	for _, version := range strings.Split(rawVersions, " ") {
		version = strings.TrimSpace(version)
		if len(version) > 0 {
			versions = append(versions, version)
		}
	}
	return versions
}

// FindVersionsInEnv returns the version from the environment if present
func FindVersionsInEnv(plugin string) ([]string, bool) {
	envVariableName := "ASDF_" + strings.ToUpper(plugin) + "_VERSION"
	versionString := os.Getenv(envVariableName)
	if versionString == "" {
		return nil, false
	}
	return ParseVersion(versionString), true
}

// FindVersionsInDir returns the version from the current directory
func FindVersionsInDir(dir string, plugin string, config Config) ([]string, bool, error) {
	filepath := path.Join(dir, toolVersionsFile)
	if _, err := os.Stat(filepath); err != nil {
		return nil, false, nil
	}
	return FindVersionsInToolFile(filepath, plugin)
}

// FindVersionsInToolFileContent returns the version of a plugin from the toolsfile content
func FindVersionsInToolFileContent(plugin string, content string) ([]string, bool) {
	for _, line := range ReadLines(content) {
		tokens := strings.SplitN(line, " ", 2)
		if strings.TrimSpace(tokens[0]) == plugin {
			return ParseVersion(tokens[1]), true
		}
	}
	return nil, false
}

// FindVersionsInToolFile returns the version of a plugin from the toolsfile at the given path
func FindVersionsInToolFile(filepath string, plugin string) ([]string, bool, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, false, err
	}
	versions, found := FindVersionsInToolFileContent(plugin, string(content))
	return versions, found, nil
}

// FindVersions returns the current versions for the plugin
func FindVersions(plugin string, config Config) ([]string, bool, error) {
	version, found := FindVersionsInEnv(plugin)
	dir, err := os.Getwd()
	if err != nil {
		return nil, false, err
	}
	for !found {
		version, found, err = FindVersionsInDir(dir, plugin, config)
		if err != nil {
			return nil, false, err
		}
		nextDir := path.Dir(dir)
		if nextDir == dir {
			break
		}
		dir = nextDir
	}
	if !found {
		homeToolsFile := path.Join(os.Getenv("HOME"), toolVersionsFile)
		if _, err := os.Stat(homeToolsFile); err == nil {
			version, found, err = FindVersionsInToolFile(homeToolsFile, plugin)
		}
	}
	return version, found, err
}

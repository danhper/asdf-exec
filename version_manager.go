package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	toolVersionsFile      string = ".tool-versions"
	legacyFileNamesScript string = "list-legacy-filenames"
	parseLegacyFileScript string = "parse-legacy-file"
)

var legacyFileNamesCache = make(map[string][]string)

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

// GetLegacyFilenames retrieves the legacy filenames of the plugin
func GetLegacyFilenames(plugin string) ([]string, error) {
	pluginPath := GetPluginPath(plugin)
	listFilenames := path.Join(pluginPath, "bin", legacyFileNamesScript)
	if _, err := os.Stat(listFilenames); err != nil {
		return []string{}, nil
	}
	result, err := exec.Command("bash", listFilenames).Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(result), " "), nil
}

// GetVersionsFromLegacyfile returns the versions in the legacy file
func GetVersionsFromLegacyfile(plugin string, legacyFilepath string) (versions []string, err error) {
	pluginPath := GetPluginPath(plugin)
	parseScriptPath := path.Join(pluginPath, "bin", parseLegacyFileScript)

	var rawVersions []byte
	if _, err := os.Stat(parseScriptPath); err == nil {
		rawVersions, err = exec.Command("bash", parseScriptPath, legacyFilepath).Output()
	} else {
		rawVersions, err = ioutil.ReadFile(legacyFilepath)
	}
	if err != nil {
		return nil, err
	}

	for _, version := range strings.Split(string(rawVersions), " ") {
		versions = append(versions, strings.TrimSpace(version))
	}
	return
}

// FindVersionsInLegacyFile returns the version from the legacy file if found
func FindVersionsInLegacyFile(dir string, plugin string) (versions []string, found bool, err error) {
	var legacyFileNames []string
	if names, ok := legacyFileNamesCache[plugin]; ok {
		legacyFileNames = names
	} else {
		legacyFileNames, err = GetLegacyFilenames(plugin)
		if err != nil {
			return
		}
		legacyFileNamesCache[plugin] = legacyFileNames
	}
	for _, filename := range legacyFileNames {
		filename = strings.TrimSpace(filename)
		filepath := path.Join(dir, filename)
		if _, err := os.Stat(filepath); err == nil {
			versions, err := GetVersionsFromLegacyfile(plugin, filepath)
			if len(versions) == 0 || (len(versions) == 1 && versions[0] == "") {
				return nil, false, nil
			}
			return versions, err == nil, err
		}
	}
	return nil, false, nil
}

// FindVersionsInDir returns the version from the current directory
func FindVersionsInDir(dir string, plugin string, config Config) ([]string, bool, error) {
	filepath := path.Join(dir, toolVersionsFile)
	if _, err := os.Stat(filepath); err == nil {
		versions, found, err := FindVersionsInToolFile(filepath, plugin)
		if found || err != nil {
			return versions, found, err
		}
	}
	if config.LegacyVersionFile {
		return FindVersionsInLegacyFile(dir, plugin)

	}
	return nil, false, nil
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

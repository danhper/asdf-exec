package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Config contains asdf configuration
type Config struct {
	// Whether to use legacy version files such as .python-version
	LegacyVersionFile bool
}

// ParseBool parses a bool from a string
func ParseBool(boolStr string) (bool, error) {
	boolStr = strings.ToLower(boolStr)
	if boolStr == "yes" || boolStr == "1" || boolStr == "true" {
		return true, nil
	} else if boolStr == "no" || boolStr == "0" || boolStr == "false" {
		return false, nil
	} else {
		return false, fmt.Errorf("unexepected boolean value: %s", boolStr)
	}
}

// Returns a config from a string
func ConfigFromString(content string) (config Config, err error) {
	for _, line := range ReadLines(content) {
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			return config, fmt.Errorf("invalid line in configuration file: ")
		}
		key := strings.TrimSpace(tokens[0])
		rawValue := strings.TrimSpace(tokens[1])
		value, err := ParseBool(rawValue)
		if err != nil {
			return config, err
		}
		if key == "legacy_version_file" {
			config.LegacyVersionFile = value
		}
	}
	return
}

// Returns a config from the given configuration file
func ConfigFromFile(filepath string) (Config, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Config{LegacyVersionFile: false}, err
	}
	return ConfigFromString(string(content))
}

// Returns a config from the default configuration file
func ConfigFromDefaultFile() (Config, error) {
	filepath := path.Join(os.Getenv("HOME"), ".asdfrc")
	return ConfigFromFile(filepath)
}

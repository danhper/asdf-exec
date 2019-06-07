package main

import (
	"path"
	"strings"
)

// ReadLines reads all the lines in a given file
// removing spaces and comments which are marked by '#'
func ReadLines(content string) (lines []string) {
	for _, line := range strings.Split(content, "\n") {
		line = strings.SplitN(line, "#", 2)[0]
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return
}

// RemoveAsdfPath returns the PATH without asdf shims path
func RemoveAsdfPath(currentPath string) string {
	paths := strings.Split(currentPath, ":")
	asdfShimPath := path.Join(GetAsdfDataPath(), "shims")
	var newPaths []string
	for _, fspath := range paths {
		if fspath != asdfShimPath {
			newPaths = append(newPaths, fspath)
		}
	}
	return strings.Join(newPaths, ":")
}

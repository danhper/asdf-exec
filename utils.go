package main

import "strings"

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

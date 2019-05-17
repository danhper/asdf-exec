package main

import (
	"fmt"
	"os"
	"syscall"
)

func findExecutablePath() (string, error) {
	config, err := ConfigFromDefaultFile()
	if err != nil {
		return "", err
	}

	shim := os.Args[1]
	shimPath := GetShimPath(shim)
	executable, found, err := FindExecutable(shimPath, config)
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("%s not found", shim)
	}
	return GetExecutablePath(executable)
}

func main() {
	executable, err := findExecutablePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: "+err.Error())
		os.Exit(1)
	}

	args := []string{executable}
	args = append(args, os.Args[2:]...)
	syscall.Exec(executable, args, os.Environ())
}

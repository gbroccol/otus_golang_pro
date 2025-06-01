package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: go-envdir <envdir> <command> [args...]")
		os.Exit(1)
	}

	envDir := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read envdir:", err)
		os.Exit(1)
	}

	code := RunCmd(command, env)
	os.Exit(code)
}

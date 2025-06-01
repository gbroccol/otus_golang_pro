package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.Contains(name, "=") {
			return nil, errors.New("invalid env var name with '='")
		}

		fullPath := filepath.Join(dir, name)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			env[name] = EnvValue{NeedRemove: true}
			continue
		}

		// take first line (up to \n or whole file)
		line := data
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			line = data[:i]
		}

		// replace terminal nulls with newlines
		line = bytes.ReplaceAll(line, []byte{0}, []byte{'\n'})

		// trim trailing spaces and tabs
		clean := strings.TrimRight(string(line), " \t")

		env[name] = EnvValue{Value: clean}
	}

	return env, nil
}

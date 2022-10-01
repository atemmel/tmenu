package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func findExecutablesInPath() []string {
	pathStr := os.Getenv("PATH")
	paths := strings.Split(pathStr, ";")

	executables := make([]string, 0, 32)
	for _, path := range paths {
		execs, err := findExecutablesInDir(path)
		if err != nil {
			panic(err)
		}

		executables = append(executables, execs...)
	}
	return executables
}

func findExecutablesInDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	found := make([]string, 0, 16)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !isExecutable(file) {
			continue
		}

		found = append(found, file.Name())
	}

	return found, nil
}

func isExecutable(file fs.FileInfo) bool {
	return file.Mode() & 0111 != 0 || strings.HasSuffix(file.Name(), ".exe")
}

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/atemmel/tmenu/pkg/tmenu"
)

const runHistoryFile = "tmenu_run_recent.json"

func Run(t *tmenu.Tmenu) {
	t.Prompt = "run"
	options := findRunnableThings()

	history, options := LookupHistory(options, runHistoryFile)

	selection := t.Repl(options)

	if selection == nil {
		return
	}

	AppendHistory(history, *selection, runHistoryFile)

	runSelected(*selection);
}

func findExecutablesInPath() []string {
	pathStr := os.Getenv("PATH")
	paths := strings.Split(pathStr, pathVarSep)

	executables := make([]string, 0, 32)
	for _, path := range paths {
		execs := findExecutablesInDir(path)
		executables = append(executables, execs...)
	}
	return executables
}

func findExecutablesInDir(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}
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

	return found
}

func runProg(prog string) {
	cmd := exec.Command(prog)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

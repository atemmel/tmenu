package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func Open(t *Tmenu) {
    dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	t.Prompt = "open project"
	options := findProjects(dir)
	selection := t.Repl(options)

	if selection == nil {
		return
	}

	*selection = dir + "/" + *selection + "/"
	*selection = strings.ReplaceAll(*selection, "\\", "/")
	fmt.Println(*selection)
	executeProjectCommand(*selection)
}

func executeProjectCommand(dir string) {
	command := "cmd /C open-project.bat"

	args := strings.Split(command, " ")
	cmd := args[0]
	args = append(args[1:], dir)

	exe := exec.Command(cmd, args...)
	err := exe.Run()
	if err != nil {
		panic(err)
	}
}

func findProjects(dir string) []string {
	projects := findProjectsInDir(dir)

	for i := range projects {
		projects[i] = projects[i][len(dir)+1:]
	}

	return projects
}

func findProjectsInDir(dir string) []string {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	projects := make([]string, 0, 16)

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == ".git" {
			return append(projects, dir)
		}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirProjects := findProjectsInDir(dir + "/" + entry.Name())
			projects = append(projects, dirProjects...)
		}
	}

	return projects
}

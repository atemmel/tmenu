package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/atemmel/tmenu/pkg/tmenu"
)

const openHistoryFile = "tmenu_open_recent.json"

func Open(t *tmenu.Tmenu) {
    dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	t.Prompt = "open project"
	options := findProjects(dir)

	var history History
	cacheDir, err := os.UserCacheDir()
	if err == nil {
		history, err = ReadHistory(cacheDir + "/" + openHistoryFile)
		if err == nil {
			in, out := HistorySplit(history, options)
			SortEntriesByHistory(history, in)
			options = append(in, out...)
		} else {
			//TODO: show error message
		}
	} else {
		//TODO: show another error message
	}

	if history == nil {
		history = make(History)
	}

	selection := t.Repl(options)

	if selection == nil {
		return
	}

	//TODO: this branch can mayyybe be avoided
	count, ok := history[*selection]
	if !ok {
		history[*selection] = 1
	} else {
		history[*selection] = count + 1
	}

	*selection = dir + "/" + *selection + "/"
	*selection = strings.ReplaceAll(*selection, "\\", "/")
	//fmt.Println(*selection)

	err = WriteHistory(history, cacheDir + "/" + openHistoryFile)
	if err != nil {
		//TODO: show yet another error message
	}

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

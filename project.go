package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/atemmel/tmenu/pkg/tmenu"
)

const projectHistoryFile = "tmenu_project_recent.json"

func Project(t *tmenu.Tmenu, dir, command string) {
	var err error
	if dir == "" {
		dir, err = os.UserHomeDir()
		if err != nil {
			panic(err)
		}
	}

	t.Prompt = "open project"
	options := findProjects(dir)

	history, options := LookupHistory(options, projectHistoryFile)

	selection := t.Repl(options)

	if selection == nil {
		return
	}

	AppendHistory(history, *selection, projectHistoryFile)

	*selection = dir + "/" + *selection + "/"
	*selection = strings.ReplaceAll(*selection, "\\", "/")

	executeProjectCommand(*selection, command)
}

func executeProjectCommand(dir, command string) {
	//command := "cmd /C open-project.bat"

	args := strings.Split(command, " ")
	cmd := args[0]
	args = append(args[1:], dir)
	fmt.Println(cmd, args)

	exe := exec.Command(cmd, args...)
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr
	err := exe.Run()
	if err != nil {
		panic(err)
	}
}

func findProjects(dir string) []string {
	t0 := time.Now()
	projects := findProjectsRecursionBase(dir)
	t1 := time.Now()

	for i := range projects {
		projects[i] = projects[i][len(dir)+1:]
	}

	fmt.Println("time spent finding repositories:", t1.Sub(t0))

	return projects
}

func findProjectsRecursionBase(base string) []string {
	entries, err := ioutil.ReadDir(base)
	if err != nil {
		panic(err)
	}

	directories := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}

	projects := []string{}
	projectsMutex := sync.Mutex{}
	group := sync.WaitGroup{}

	group.Add(len(directories))

	find := func(dir string) {
		defer group.Done()
		newProjects := findProjectsRecursion(dir)
		projectsMutex.Lock()
		defer projectsMutex.Unlock()
		projects = append(projects, newProjects...)
	}

	for _, dir := range directories {
		go find(base + "/" + dir)
	}
	group.Wait()

	return projects
}

func findProjectsRecursion(dir string) []string {
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
			dirProjects := findProjectsRecursion(dir + "/" + entry.Name())
			projects = append(projects, dirProjects...)
		}
	}

	return projects
}

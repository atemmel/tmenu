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

const openHistoryFile = "tmenu_open_recent.json"

func Project(t *tmenu.Tmenu, dir string) {
	var err error
	if dir == "" {
		dir, err = os.UserHomeDir()
		if err != nil {
			panic(err)
		}
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

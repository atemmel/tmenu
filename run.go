package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var shortcutStr = "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs"

func Run(t *Tmenu) {
	options := findRunnableThings()
	selection := t.Repl(options)

	if selection == nil {
		return
	}
	fmt.Println(*selection)

	if strings.HasSuffix(*selection, ".lnk") {
		runLnk(*selection)
	} else {
		runProg(*selection)
	}
}

func runLnk(lnk string) {
	lnk = shortcutStr + "\\" + lnk
	fmt.Println("running lnk", lnk)
	cmd := exec.Command("cmd.exe", "start", lnk)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func runProg(prog string) {
	fmt.Println("running prog", prog)
	cmd := exec.Command(prog)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func findRunnableThings() []string {
	execs := findExecutablesInPath()
	shortcuts := findShortcuts()
	return append(execs, shortcuts...)
}

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

func findShortcuts() []string {
	shortcuts := findShortcutsInDir(shortcutStr)
	for i, s := range shortcuts {
		shortcuts[i] = s[len(shortcutStr) + 1:]
	}
	return shortcuts
}

func findShortcutsInDir(dir string) []string {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	shortcuts := make([]string, 0, 32)

	for _, entry := range entries {
		if entry.IsDir() {
			dirShortcuts := findShortcutsInDir(dir + "\\" + entry.Name())
			shortcuts = append(shortcuts, dirShortcuts...)
		} else {
			shortcuts = append(shortcuts, dir + "\\" + entry.Name())
		}
	}

	return shortcuts
}

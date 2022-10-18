package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/atemmel/tmenu/pkg/tmenu"
	"golang.org/x/sys/windows"
)

var shortcutDirs = []string{
	"C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs",
}
const runHistoryFile = "tmenu_run_recent.json"

func Run(t *tmenu.Tmenu) {
	cacheDir, err := os.UserConfigDir()
	if err == nil {
		startmenuFolder2 := cacheDir + "\\Microsoft\\Windows\\Start Menu\\Programs"
		shortcutDirs = append(shortcutDirs, startmenuFolder2)
	}

	t.Prompt = "run"
	options, dirMap := findRunnableThings()

	history, options := LookupHistory(options, runHistoryFile)
	_ = history

	selection := t.Repl(options)

	if selection == nil {
		return
	}

	AppendHistory(history, *selection, runHistoryFile)

	if strings.HasSuffix(*selection, ".lnk") {
		idx := dirMap[*selection]
		runLnk(*selection, idx)
	} else {
		runProg(*selection)
	}
}

func runLnk(lnk string, idx int) {
	base := shortcutDirs[idx]
	lnk = base + "\\" + lnk
	err := shellExecute(lnk)
	if err != nil {
		panic(err)
	}
}

func shellExecute(lnk string) error {
	var handle windows.Handle = 0;
	args := windows.StringToUTF16Ptr(lnk)
	var show int32 = 9
	return windows.ShellExecute(handle, nil, args, nil, nil, show)
}

func runProg(prog string) {
	cmd := exec.Command(prog)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func findRunnableThings() ([]string, map[string]int) {
	execs := findExecutablesInPath()
	shortcuts, dirMap := findShortcuts()
	return append(execs, shortcuts...), dirMap
}

func findExecutablesInPath() []string {
	pathStr := os.Getenv("PATH")
	paths := strings.Split(pathStr, ";")

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

func findShortcuts() ([]string, map[string]int) {
	dirMap := make(map[string]int, 128);
	shortcuts := make([]string, 0, 128)
	for dirIdx, dir := range shortcutDirs {
		newShortcuts := findShortcutsInDir(dir)
		for i, s := range newShortcuts {
			r := s[len(dir) + 1:]
			newShortcuts[i] = r
			dirMap[r] = dirIdx
		}
		shortcuts = append(shortcuts, newShortcuts...)
	}
	return shortcuts, dirMap
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

func isExecutable(file fs.FileInfo) bool {
	return file.Mode() & 0111 != 0 || strings.HasSuffix(file.Name(), ".exe")
}

// +build windows

package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

const pathVarSep = ";"

var shortcutDirs = []string{
	"C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs",
}

var dirMap = make(map[string]int)

func runSelected(selection string) {
	if strings.HasSuffix(selection, ".lnk") {
		idx := dirMap[selection]
		runLnk(selection, idx)
	} else {
		runProg(selection)
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

func findRunnableThings() []string {
	cacheDir, err := os.UserConfigDir()
	if err == nil {
		startmenuFolder2 := cacheDir + "\\Microsoft\\Windows\\Start Menu\\Programs"
		shortcutDirs = append(shortcutDirs, startmenuFolder2)
	}
	execs := findExecutablesInPath()
	shortcuts := []string{}
	shortcuts, dirMap = findShortcuts()
	return append(execs, shortcuts...)
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

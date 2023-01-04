// !build windows

package main

import (
	"io/fs"
)

const pathVarSep = ":"

func findRunnableThings() []string {
	return findExecutablesInPath()
}

func runSelected(selection string) {
	runProg(selection)
}

func isExecutable(file fs.FileInfo) bool {
	return file.Mode() & 0111 != 0
}

package main

import (
	"fmt"
)

func Open(t *Tmenu) {
	t.Prompt = "open project"
	options := findProjects()
	selection := t.Repl(options)

	if selection == nil {
		return
	}
	fmt.Println(*selection)
}

func findProjects() []string {
	return []string{}
}

package main

import (
	"fmt"

	"github.com/atemmel/tmenu/pkg/tmenu"
)

func Stdin(t *tmenu.Tmenu, options []string) {
	selection := t.Repl(options)
	if selection == nil {
		return
	}
	fmt.Println(*selection)
}

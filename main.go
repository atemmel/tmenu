package main

import (
	_ "embed"
	"flag"
	"fmt"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	defaultWidth = 80
	defaultHeight = 10
)

//go:embed fonts/FiraCodeNerdFont-regular.ttf
var defaultFontBytes []byte

func init() {
	flag.Parse()
}

func main() {
	runtime.LockOSThread()
	//TODO: only init what is needed
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	rw, err := sdl.RWFromMem(defaultFontBytes)
	if err != nil {
		panic(err)
	}

	font, err := ttf.OpenFontRW(rw, 1, 16)
	if err != nil {
		panic(err)
	}

	tmenu, err := NewTmenu(defaultWidth, defaultHeight, font)
	if err != nil {
		panic(err)
	}
	defer tmenu.Destroy()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("No command provided")
		return
	}

	cmd := args[0]
	if cmd == "run" {
		Run(tmenu)
	} else if cmd == "open" {
		Open(tmenu)
	} else {
		fmt.Println("Unrecognized command,", cmd)
	}

}

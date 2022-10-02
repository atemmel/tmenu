package main

import (
	_ "embed"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	defaultWidth = 800
	defaultHeight = 600
)

//go:embed fonts/FiraCodeNerdFont-regular.ttf
var defaultFontBytes []byte

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

	Run(tmenu)
}

package main

import (
	_ "embed"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"runtime"
	"time"
)

const (
	defaultWidth = 600
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

	font, err := ttf.OpenFontRW(rw, 1, 12)
	if err != nil {
		panic(err)
	}

	tmenu, err := NewTmenu(defaultWidth, defaultHeight, font)
	if err != nil {
		panic(err)
	}
	defer tmenu.Destroy()

	tmenu.Redraw()
	for tmenu.IsRunning() {
		time.Sleep(167)
		if !tmenu.PollEvents() {
			continue
		}
		tmenu.Redraw()
	}
}

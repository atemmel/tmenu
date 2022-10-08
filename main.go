package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/atemmel/tmenu/pkg/tmenu"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	defaultWidth = 80
	defaultHeight = 10
)

type Mode int

const (
	RunMode Mode = iota
	OpenMode
	NoneProvidedMode
	UnknownMode
)

var (
//go:embed fonts/FiraCodeNerdFont-regular.ttf
	defaultFontBytes []byte
	currentMode Mode = UnknownMode
)

func init() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		currentMode = NoneProvidedMode
		return
	}
	cmd := args[0]
	switch cmd {
	case "run":
		currentMode = RunMode
	case "open":
		currentMode = OpenMode
	default:
		currentMode = UnknownMode
	}
}

func main() {
	if currentMode == NoneProvidedMode {
		fmt.Println("No command provided")
		os.Exit(1)
	}
	if currentMode == UnknownMode {
		fmt.Println("Unknown mode:", flag.Args()[0])
		os.Exit(2)
	}

	tmenu := initSdlAndTmenu()
	defer sdl.Quit()
	defer tmenu.Destroy()

	switch currentMode {
		case RunMode:
			Run(tmenu)
		case OpenMode:
			Open(tmenu)
	}
}

func initSdlAndTmenu() *tmenu.Tmenu {
	runtime.LockOSThread()

	//TODO: only init what is needed
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

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

	tmenu, err := tmenu.NewTmenu(defaultWidth, defaultHeight, font)
	if err != nil {
		panic(err)
	}
	return tmenu
}


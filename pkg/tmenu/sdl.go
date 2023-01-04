package tmenu

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)


func initSdl() {
	runtime.LockOSThread()

	if err := sdl.Init(sdlFlags); err != nil {
		panic(err)
	}

	if err := ttf.Init(); err != nil {
		panic(err)
	}
}

func loadFont() *ttf.Font {
	rw, err := sdl.RWFromMem(defaultFontBytes)
	if err != nil {
		panic(err)
	}

	font, err := ttf.OpenFontRW(rw, 1, 16)
	if err != nil {
		panic(err)
	}

	return font
}

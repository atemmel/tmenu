package main

import (
	"bufio"
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
	StdinMode
	NoneProvidedMode
	UnknownMode
)

var (
//go:embed fonts/FiraCodeNerdFont-regular.ttf
	defaultFontBytes []byte
	stdin []string
	currentMode Mode = UnknownMode
	promptOverride string
)

func readStdin() []string {
	file, err := os.Stdin.Stat()
	if err != nil {
		//TODO: investigate why this is
		// this is a weird windows-only fix btw
		return []string{}
	}
	if file.Mode() & os.ModeNamedPipe == 0 {
		return []string{}
	}

	reader := bufio.NewReader(os.Stdin)
	lines := []string{}
	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		if len(line) > 0 && line[len(line) - 1] == '\n' {
			line = line[:len(line)-1]
		} else {
			continue
		}
		if len(line) == 0 {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func init() {
	flag.StringVar(&promptOverride, "prompt", "", "Override default prompt title")
	flag.Parse()
	stdin = readStdin()
	if len(stdin) > 0 {
		currentMode = StdinMode
		return
	}

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
		fmt.Fprintln(os.Stderr, "No command provided")
		os.Exit(1)
	}
	if currentMode == UnknownMode {
		fmt.Fprintln(os.Stderr, "Unknown mode:", flag.Args()[0])
		os.Exit(2)
	}

	tmenu := initSdlAndTmenu()
	defer sdl.Quit()
	defer tmenu.Destroy()

	if promptOverride != "" {
		tmenu.Prompt = promptOverride
	}

	switch currentMode {
		case RunMode:
			Run(tmenu)
		case OpenMode:
			Open(tmenu)
		case StdinMode:
			Stdin(tmenu, stdin)
	}
}

func initSdlAndTmenu() *tmenu.Tmenu {
	runtime.LockOSThread()

	const sdlFlags uint32 = sdl.INIT_EVENTS | sdl.INIT_VIDEO;

	if err := sdl.Init(sdlFlags); err != nil {
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


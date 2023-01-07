package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

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
	ProjectMode
	StdinMode
	NoneProvidedMode
	UnknownMode
)

var (
	stdin []string
	currentMode Mode = UnknownMode
	promptOverride string
	dirOverride string
	command string
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
	flag.StringVar(&dirOverride, "dir", "", "Override default operating directory")
	flag.StringVar(&command, "command", "", "Command to run after selecting a project")
	flag.BoolVar(&tmenu.DebugPerformance, "perf", false, "Debug performance (quit after a single window loop iteration)")
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
	case "project":
		currentMode = ProjectMode
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

	tmenu := tmenu.NewTmenu(defaultWidth, defaultHeight)
	defer sdl.Quit()
	defer ttf.Quit()
	defer tmenu.Destroy()

	if promptOverride != "" {
		tmenu.Prompt = promptOverride
	}

	switch currentMode {
		case RunMode:
			Run(&tmenu)
		case ProjectMode:
			Project(&tmenu, dirOverride, command)
		case StdinMode:
			Stdin(&tmenu, stdin)
	}
}

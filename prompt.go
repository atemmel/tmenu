package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	paddingX = 12
	paddingY = 8
)

var (
	textColor = sdl.Color{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	backgroundColor = sdl.Color{
		R: 66,
		G: 66,
		B: 82,
		A: 255,
	}

	selectedBackgroundColor = sdl.Color{
		R: 106,
		G: 106,
		B: 122,
		A: 255,
	}
)

type Tmenu struct {
	running bool
	submitted bool
	window *sdl.Window
	renderer *sdl.Renderer
	font *ttf.Font

	prompt string
	input string
	options []string
	filteredOptions []string
	selectedIndex int
	w int32
	h int32
	textH int32
}

func NewTmenu(w, h int, font *ttf.Font, options []string) (*Tmenu, error) {
	dm, err := sdl.GetDesktopDisplayMode(0)
	if err != nil {
		return nil, err
	}

	x := dm.W / 2 - int32(w) / 2
	var y int32 = 100

	window, err := sdl.CreateWindow(
		"tmenu", 
		x, 
		y, 
		int32(w), 
		int32(h), 
		sdl.WINDOW_SHOWN | sdl.WINDOW_BORDERLESS)
	if err != nil {
		return nil, err
	}

	window.Raise()
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		return nil, err
	}
	sdl.StartTextInput()

	filteredOptions := make([]string, len(options))
	copy(filteredOptions, options)

	return &Tmenu{
		running: true,
		submitted: false,
		window: window,
		renderer: renderer,
		font: font,
		prompt: "prompt",
		input: "",
		w: int32(w),
		h: int32(h),
		textH: int32(font.Ascent() - font.Descent()),
		options: options,
		filteredOptions: filteredOptions,
		selectedIndex: 0,
	}, nil
}

func (t *Tmenu) Destroy() {
	t.renderer.Destroy()
	t.window.Destroy()
}

func (t *Tmenu) IsRunning() bool {
	return t.running
}

func (t *Tmenu) drawBackground(x, y int32, clr sdl.Color) {
	inputBackground := sdl.Rect{
		X: x,
		Y: y,
		W: t.w,
		H: t.textH + paddingY / 2,
	}

	t.renderer.SetDrawColor(clr.R, clr.G, clr.B, clr.A)
	t.renderer.FillRect(&inputBackground)
}

func (t *Tmenu) drawText(text string, x, y int32) {
	textSurface, err := t.font.RenderUTF8Blended(text, textColor)
	if err != nil {
		panic(err)
	}
	defer textSurface.Free()
	texture, err := t.renderer.CreateTextureFromSurface(textSurface)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()
	textSize := textSurface.Bounds().Size()
	textRect := sdl.Rect{
		X: x,
		Y: y,
		W: int32(textSize.X),
		H: int32(textSize.Y),
	}

	t.renderer.Copy(texture, nil, &textRect)
}

func (t *Tmenu) drawPrompt() {
	str := fmt.Sprintf("%s: %s", t.prompt, t.input)
	t.drawRow(0, str, nil)
}

func (t *Tmenu) drawOptions() {
	for i, opt := range t.filteredOptions {
		var clr *sdl.Color = nil
		if i == t.selectedIndex {
			clr = &selectedBackgroundColor
		} 
		t.drawRow(int32(i + 1), opt, clr)
	}
}

func (t *Tmenu) drawRow(i int32, content string, clr *sdl.Color) {
	inputY0 := paddingY + int32(paddingY / 2)
	inputY1 := inputY0 + (t.textH + paddingY / 2) * i

	if clr != nil {
		t.drawBackground(0, inputY1 - paddingY / 2, *clr)
	}
	t.drawText(content, paddingX, inputY1 - paddingY / 4)
}

func (t *Tmenu) GetSelection() *string {
	if len(t.filteredOptions) == 0 || !t.submitted {
		return nil
	}
	return &t.filteredOptions[t.selectedIndex]
}

func (t *Tmenu) Redraw() {
	t.renderer.SetDrawColor(backgroundColor.R, 
		backgroundColor.G, 
		backgroundColor.B,
		backgroundColor.A)
	t.renderer.Clear()

	t.drawPrompt()
	t.drawOptions()

	t.renderer.Present()
}

func (t *Tmenu) quit() {
	t.running = false
}

func (t *Tmenu) insertInput(input string) {
	t.input += input
	t.refilter()
}

func (t *Tmenu) refilter() {
	t.filteredOptions = filter(t.input, t.options)
}

func (t *Tmenu) erase() {
	if len(t.input) == 0 {
		return
	}
	t.input = t.input[0:len(t.input) - 1]
	t.refilter()
}

func (t *Tmenu) submit() {
	t.submitted = true
	t.quit()
}

func (t *Tmenu) PollEvents() (updated bool) {
	updated = false
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			t.quit()
		case *sdl.KeyboardEvent:
			keyEvent := event.(*sdl.KeyboardEvent)
			if keyEvent.Type == sdl.KEYDOWN && keyEvent.Repeat == 0 {
				fmt.Println("KeyboardEvent:", keyEvent)
				t.handleKeys(keyEvent)
				updated = true
			}
		case *sdl.TextInputEvent:
			textInputEvent := event.(*sdl.TextInputEvent)
			t.insertInput(textInputEvent.GetText())
			updated = true
		case *sdl.TextEditingEvent:
			fmt.Println("TextEditingEvent:", event.(*sdl.TextEditingEvent))
			updated = true
		}
	}
	return updated
}

func (t *Tmenu) moveCursorUp() {
	if t.selectedIndex > 0 {
		t.selectedIndex--
	} else if len(t.options) > 0 {
		t.selectedIndex = len(t.options) - 1
	} else {
		t.selectedIndex = 0
	}
}

func (t *Tmenu) moveCursorDown() {
	if t.selectedIndex < len(t.options) - 1 {
		t.selectedIndex++
	} else {
		t.selectedIndex = 0
	}
}

func (t *Tmenu) handleKeys(key *sdl.KeyboardEvent) {
	switch key.Keysym.Sym {
	case sdl.K_ESCAPE:
		t.quit()
	case sdl.K_TAB:
		shiftDown := key.Keysym.Mod & sdl.KMOD_SHIFT != 0
		if shiftDown {
			t.moveCursorUp()
		} else {
			t.moveCursorDown()
		}
	case sdl.K_UP:
		t.moveCursorUp()
	case sdl.K_DOWN:
		t.moveCursorDown()
	case sdl.K_BACKSPACE:
		t.erase()
	case sdl.K_RETURN:
		t.submit()
	}
}

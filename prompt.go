package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Tmenu struct {
	running bool
	window *sdl.Window
	font *ttf.Font
}

func NewTmenu(w, h int, font *ttf.Font) (*Tmenu, error) {
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
	sdl.StartTextInput()


	return &Tmenu{
		window: window,
		running: true,
		font: font,
	}, nil
}

func (t *Tmenu) Destroy() {
	t.window.Destroy()
}

func (t *Tmenu) IsRunning() bool {
	return t.running
}

func (t *Tmenu) Redraw() {
	//TODO: this

	surface, err := t.window.GetSurface()
	if err != nil {
		t.window.Destroy()
		t.running = false
		return
	}
	surface.SetBlendMode(sdl.BLENDMODE_NONE)
	surface.FillRect(nil, 0xff444455)

	clr := sdl.Color{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	str := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed egestas mi nec euismod iaculis. Nullam suscipit massa tempor maximus ornare. Sed placerat rhoncus rhoncus. Duis dui purus, commodo ac purus quis, facilisis vehicula ipsum. Aenean mi justo, iaculis ac dapibus in, vestibulum id augue. Sed auctor viverra arcu nec scelerisque. Pellentesque viverra eros ut massa condimentum varius. Etiam quis condimentum metus. Aenean ultrices turpis quis turpis bibendum aliquam. Quisque ut eleifend leo. Sed at nisl at lorem laoreet tincidunt vel non lorem. Suspendisse in turpis quam. Sed aliquet, odio ut eleifend suscipit, risus nisl mollis lorem, et sollicitudin mi tellus at orci. Donec rhoncus mauris ac vulputate lobortis. Vivamus id scelerisque ipsum. Duis ipsum nulla, tincidunt id finibus in, convallis sit amet lorem. 
`

	text, err := t.font.RenderUTF8Solid(str, clr)
	if err != nil {
		panic(err)
	}
	textSize := text.Bounds().Size()
	textRect := sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(textSize.X),
		H: int32(textSize.Y),
	}
	surfaceBounds := surface.Bounds().Size()
	surfaceRect := sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(surfaceBounds.X),
		H: int32(surfaceBounds.Y),
	}
	text.Blit(&textRect, surface, &surfaceRect)

	t.window.UpdateSurface()
}

func (t *Tmenu) quit() {
	t.running = false
}

func (t *Tmenu) PollEvents() (updated bool) {
	updated = false
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			t.quit()
			updated = true
		case *sdl.KeyboardEvent:
			keyEvent := event.(*sdl.KeyboardEvent)
			fmt.Println("KeyboardEvent:", keyEvent)
			t.handleKeys(keyEvent)
			updated = true
		case *sdl.TextInputEvent:
			fmt.Println("TextInputEvent:", event.(*sdl.TextInputEvent))
			updated = true
		case *sdl.TextEditingEvent:
			fmt.Println("TextEditingEvent:", event.(*sdl.TextEditingEvent))
			updated = true
		}
	}
	return updated
}

func (t *Tmenu) handleKeys(key *sdl.KeyboardEvent) {
	if key.Keysym.Sym == sdl.K_ESCAPE {
		t.quit()
	}
}

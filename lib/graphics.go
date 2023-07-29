package lib

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	BLOCK = reflect.TypeOf(Color(0)).Size()
)

type Color = uint32

type Graphics struct {
	Renderer *sdl.Renderer
	Texture  *sdl.Texture
	Window   *sdl.Window
	Width    int32
	Height   int32

	pixels []Color
}

func (g Graphics) Render() {
	g.Renderer.Clear()
	g.Texture.Update(nil, unsafe.Pointer(&g.pixels[0]), int(BLOCK * uintptr(g.Width)))
	g.Renderer.Copy(g.Texture, nil, nil)
	g.Renderer.Present()
}

func (g Graphics) Clear(color Color) {
	for i := 0; i < len(g.pixels); i++ {
		g.pixels[i] = color
	}
}

func (g Graphics) Destroy() {
	g.Texture.Destroy()
	g.Renderer.Destroy()
	g.Window.Destroy()
}

func NewGraphics(width int32, height int32) (*Graphics, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, errors.New("failed to initialize SDL")
	}

	var window *sdl.Window
	window, err := sdl.CreateWindow("3f", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width, height, 0)
	if err != nil {
		return nil, err
	}

	var renderer *sdl.Renderer
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_TARGETTEXTURE)
	if err != nil {
		window.Destroy()
		return nil, err
	}

	var texture *sdl.Texture
	texture, err = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		renderer.Destroy()
		window.Destroy()
		return nil, err
	}

	return &Graphics{
		Renderer: renderer,
		Texture:  texture,
		Window:   window,
		Width:    width,
		Height:   height,
		pixels:   make([]Color, width*height),
	}, nil
}

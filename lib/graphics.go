package lib

import (
	"errors"
	"math"
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
	g.Texture.Update(nil, unsafe.Pointer(&g.pixels[0]), int(BLOCK*uintptr(g.Width)))
	g.Renderer.Copy(g.Texture, nil, nil)
	g.Renderer.Present()
}

func (g Graphics) Clear(color Color) {
	for i := 0; i < len(g.pixels); i++ {
		g.pixels[i] = color
	}
}

func (g Graphics) PutPixel(x int32, y int32, color Color) {
	g.pixels[g.Width*y+x] = color
}

func (g Graphics) Line(x1, y1, x2, y2 int32, color Color) {
	// Adapted from https://www.ercankoclar.com/wp-content/uploads/2016/12/Bresenhams-Algorithm.pdf
	dy := math.Abs(float64(y2 - y1))
	dx := math.Abs(float64(x2 - x1))
	steep := dy > dx

	if steep {
		// we use y as the driving axis
		if y2 < y1 {
			x1, y1, x2, y2 = x2, y2, x1, y1
		}

		m := dx / dy
		e := m - 1
		x := x1

		var d int32 = 1
		if x2 < x1 {
			d = -1
		}

		for y := y1; y <= y2; y++ {
			g.PutPixel(x, y, color)
			if e >= 0 {
				x += d
				e -= 1
			}

			e += m
		}

	} else {
		if x2 < x1 {
			x1, y1, x2, y2 = x2, y2, x1, y1
		}

		m := dy / dx
		e := m - 1
		y := y1

		var d int32 = 1
		if y2 < y1 {
			d = -1
		}

		for x := x1; x <= x2; x++ {
			g.PutPixel(x, y, color)
			if e >= 0 {
				y += d
				e -= 1
			}

			e += m
		}
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

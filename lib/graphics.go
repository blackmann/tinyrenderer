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
	if x < 0 || y < 0 || x >= g.Width || y >= g.Height {
		return
	}
	g.pixels[g.Width*y+x] = color
}

func (g Graphics) Line(p1, p2 Vector3, color Color) {
	p1Flat := p1.To2()
	p2Flat := (p2.To2())

	x1 := int32(p1Flat.X)
	x2 := int32(p2Flat.X)
	y1 := int32(p1Flat.Y)
	y2 := int32(p2Flat.Y)

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

// Reference: https://www.youtube.com/watch?v=HYAgJN3x4GA
func pointInTriangle(A, B, C, point Vector2) bool {
	s1 := C.Y - A.Y
	s2 := C.X - A.X
	s3 := B.Y - A.Y
	s4 := point.Y - A.Y

	w1 := (A.X*s1 + s4*s2 - point.X*s1) / (s3*s2 - (B.X-A.X)*s1)
	w2 := (s4 - w1*s3) / s1

	return w1 >= 0 && w2 >= 0 && (w1+w2) <= 1
}

func (g Graphics) Triangle(a, b, c Vector3, color Color) {
	topLeft, bottomRight := BoundingBox(a.To2(), b.To2(), c.To2())

	for y := topLeft.Y; y < bottomRight.Y; y++ {
		for x := topLeft.X; x < bottomRight.X; x++ {
			if pointInTriangle(a.To2(), b.To2(), c.To2(), Vector2{X: x, Y: y}) {
				g.PutPixel(int32(x), int32(y), color)
			}
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

func BoundingBox(vertices ...Vector2) (topLeft, bottomRight Vector2) {
	// the top-left vertex is supposed to something like (0,0)
	// but then the triangle could be in the middle of the screen
	// so it'll be hard to derive the topmost-leftmost vertex if we start
	// from (0, 0). So we start from the bottom-right (out-of-screen)
	// and compare for the minimum
	topLeft = Vector2{X: math.Inf(0), Y: math.Inf(0)}

	for _, vertex := range vertices {
		topLeft.X = math.Min(topLeft.X, vertex.X)
		topLeft.Y = math.Min(topLeft.Y, vertex.Y)

		bottomRight.X = math.Max(bottomRight.X, vertex.X)
		bottomRight.Y = math.Max(bottomRight.Y, vertex.Y)
	}

	return
}

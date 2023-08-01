package lib

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

type App struct {
	g   *Graphics
	fps uint8

	lastRenderTime uint64
}

func (app App) Run() error {
	running := true
	for running {
		app.g.Clear(0x222222ff)
		if err := app.handleKeyInput(); err != nil {
			running = false
		}

		app.Update()

		interval := 1000 / int64(app.fps)

		tick := sdl.GetTicks64()
		// Is it possible `tick` could become more than `targetTime` because
		// some work too long?
		targetTime := app.lastRenderTime + uint64(interval)
		if tick < targetTime {
			sdl.Delay(uint32(targetTime - tick))
		}

		app.g.Render()

		app.lastRenderTime = targetTime
	}

	return nil
}

func (app App) Update() {
	app.g.Line(10, 10, 100, 80, White)
	app.g.Line(100, 500, 110, 90, Red)
	app.g.Line(300, 10, 10, 500, Green)
	app.g.Line(500, 10,100, 10,  Blue)
}

func (app App) handleKeyInput() error {
	event := sdl.PollEvent()
	if event == nil {
		return nil
	}

	switch t := event.(type) {
	case *sdl.WindowEvent:
		if t.Event == sdl.WINDOWEVENT_CLOSE {
			return errors.New("app quit")
		}
	case *sdl.KeyboardEvent:
		if t.Keysym.Sym == sdl.K_ESCAPE {
			return errors.New("app quit")
		}
	}

	return nil
}

func NewApp(g *Graphics, fps uint8) *App {
	return &App{g: g, fps: fps}
}

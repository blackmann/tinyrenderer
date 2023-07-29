package lib

import "github.com/veandco/go-sdl2/sdl"

type App struct {
	g *Graphics
}

func (app App) Run() error {
	running := true
	for running {
		event := sdl.PollEvent()
		if event == nil {
			continue
		}

		switch t := event.(type) {
		case *sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_CLOSE {
				running = false
			}
		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_ESCAPE {
				running = false
			}
		}
	}

	return nil
}

func NewApp(g *Graphics) *App {
	return &App{g}
}

package lib

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

type App struct {
	g   *Graphics
	fps uint8

	lastRenderTime uint64
	models         []*Model
}

func (app *App) init() {
	model, err := LoadModel("assets/barbarian.obj")
	if err != nil {
		app.Clean()
		panic(err)
	}

	app.models = append(app.models, model)
}

func (app App) Clean() {
	app.g.Destroy()
}

func (app App) Run() error {
	app.init()

	running := true
	defer app.Clean()

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
	translate := Vector3{X: 400, Y: 550}

	for _, model := range app.models {
		for _, object := range model.Objects {
			for _, face := range object.Faces {
				p1 := model.Vertices[face[0] - 1].Add(translate)
				p2 := model.Vertices[face[1] - 1].Add(translate)
				p3 := model.Vertices[face[2] - 1].Add(translate)
				app.g.Triangle(p1, p2, p3, White)
			}
		}
	}
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

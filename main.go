package main

import (
	"github.com/blackmann/tinyrenderer/lib"
)

// [ ] Set up SDL
// [ ] Draw line
// [ ] Implement model
func main() {
	var g *lib.Graphics
	g, err := lib.NewGraphics(800, 600)
	if err != nil {
		println("[renderer] error setting up graphics")
	}

	defer g.Destroy()
	lib.NewApp(g, 60).Run()
}

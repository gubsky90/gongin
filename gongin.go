package gongin

import (
	"time"
	"fmt"
	"github.com/gubsky90/gongin/render"
)

type Config struct {

}

type Gongin struct {
	handlers map[string][]func()
}

func New(conf Config) *Gongin {
	g := Gongin{}
	g.handlers = make(map[string][]func())
	return &g
}

func (g *Gongin) On(name string, cb func()) {
	if _, ok := g.handlers[name]; !ok {
		g.handlers[name] = make([]func(), 0)
	}
	g.handlers[name] = append(g.handlers[name], cb)
}

func (g *Gongin) fire(name string) {
	if cbs, ok := g.handlers[name]; ok {
		for _, cb := range cbs {
			cb()
		}
	}
}

func (g *Gongin) Run() {
	r := render.New(render.Config{
		Width: 640,
		Height: 480,
	})
	defer r.Destroy()

	fb := render.NewFramebuffer(640, 480)

	meshRaster := render.NewRasterizator(fb, meshShader)

	postRaster := render.NewRasterizator(r.GetWindow(), postShader)
	postRaster.SetTexture("screenTexture", fb.Color)

	mesh := render.NewMeshFromFile("teapot.obj")
	g.fire("ready")

	for !r.ShouldClose() {
		start := time.Now()
		// r.Clear()

		meshRaster.DrawMesh(mesh)
		postRaster.DrawRect()

		r.SwapBuffers()
		fmt.Printf("Frame time: %s\n", time.Since(start))
	}
}

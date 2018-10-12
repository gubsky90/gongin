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
	r := render.New()
	defer r.Destroy()

	shader := render.NewShader(shaderSource)
	meshRaster := render.NewRasterizator(shader)

	mesh := render.NewMeshFromFile("teapot.obj")

	g.fire("ready")

	for !r.ShouldClose() {
		start := time.Now()
		r.Clear()

		meshRaster.DrawMesh(mesh)

		r.SwapBuffers()
		fmt.Printf("Frame time: %s\n", time.Since(start))
	}
}

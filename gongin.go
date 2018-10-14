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

	var iTime float32

	fb := render.NewFramebuffer(640, 480)
	win := r.GetWindow()

	meshRaster := render.NewRasterizator(fb, render.NewShader(meshShader))
	meshRaster.Bind("iTime", &iTime)

	postRaster := render.NewRasterizator(win, render.NewShader(postShader))
	postRaster.Bind("screenTexture", fb.Color)
	postRaster.Bind("iTime", &iTime)

	mesh := render.NewMeshFromFile("teapot.obj")
	g.fire("ready")

	for !r.ShouldClose() {
		start := time.Now()

		iTime = float32(getTime() - getStartTime()) / 1000

		fb.Clear()
		meshRaster.Draw(mesh)

		// win.Clear()
		postRaster.DrawRect()

		r.SwapBuffers()
		fmt.Printf("Frame time: %s\n", time.Since(start))
	}
}

var startTime int64

func getStartTime() int64 {
	if startTime == 0 {
		startTime = getTime()
	}
	return startTime
}

func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
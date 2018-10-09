package gongin

import (
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

var points = []float32 {
	-0.5,  0.5,  0.0,
	-0.5, -0.5,  0.0,
	 0.5,  0.5,  0.0,

	 0.5,  0.5,  0.0,
	-0.5, -0.5,  0.0,
	 0.5, -0.5,  0.0,
}

func (g *Gongin) Run() {
	r := render.New()
	defer r.Close()

	// run := true

	// r.On("close", func(){
	// 	run = false
	// })

	mesh := render.NewMesh(points)

	shader := render.NewShader(render.ShaderSource{
		Vertex: `
			#version 410
			in vec3 vp;
			void main() {
				gl_Position = vec4(vp, 1.0);
			}
		`,
		Fragment: `
			#version 410
			uniform vec3 iColor;
			uniform float iTime;
			out vec4 frag_color;
			void main() {
				// gl_FragCoord
				vec3 col = vec3(1.0, iTime, 0.0);
				// col *= abs(sin(iTime));
				frag_color = vec4(col, 1.0);
			}
		`,
	})

	r.UseShader(shader)


	g.fire("ready")

	for !r.ShouldClose() {
		r.Draw(mesh)
	}
}

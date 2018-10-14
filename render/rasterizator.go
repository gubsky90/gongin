package render

import (
	// "os"
	// "fmt"
	// "github.com/go-gl/gl/v4.1-core/gl"
)

type Drawable interface {
	Draw()
}

// Rasterizator ...
type Rasterizator struct {
	render *Render
	shader *Shader
	target RenderTarget
}

func NewRasterizator(target RenderTarget, shader *Shader) *Rasterizator {
	r := Rasterizator{}
	r.target = target
	r.shader = shader
	return &r
}

func (r *Rasterizator) Bind(name string, value interface{}) {
	r.shader.Bind(name, value)
}

func (r *Rasterizator) before() {
	r.target.SetAsCurrentRenderTarget()
	r.shader.Use()
}

func (r *Rasterizator) after() {

}

var rect *Mesh
func (r *Rasterizator) DrawRect() {
	if rect == nil {
		rect = NewMesh([]float32{
			-1.0, -1.0, 0.0,
			 1.0, -1.0, 0.0,
			-1.0,  1.0, 0.0,
			 1.0,  1.0, 0.0,
		}, []uint32{
			0, 1, 2, 2, 1, 3,
		})
	}
	r.Draw(rect)
}

func (r *Rasterizator) Draw(drawable Drawable) {
	r.before()
	drawable.Draw()
	r.after()
}

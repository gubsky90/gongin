package render

import (
	// "os"
	// "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Drawable interface {
	Draw()
}

// Rasterizator ...
type Rasterizator struct {
	render *Render
	shader *Shader
	target RenderTarget
	textures map[string]*Texture
}

func NewRasterizator(target RenderTarget, shader *Shader) *Rasterizator {
	r := Rasterizator{}
	r.target = target
	r.shader = shader
	r.textures = make(map[string]*Texture)
	return &r
}

func (r *Rasterizator) SetTexture(name string, texture *Texture) {
	r.textures[name] = texture
}

func (r *Rasterizator) before() {
	r.shader.Use()
	r.target.SetAsCurrentRenderTarget()

	texture_slots := []uint32{
		gl.TEXTURE0,
		gl.TEXTURE1,
		gl.TEXTURE2,
	}

	var idx int32 = 0
	for name, texture := range r.textures {
		gl.ActiveTexture(texture_slots[idx])
		gl.BindTexture(gl.TEXTURE_2D, texture.GetId())
		r.shader.Set1i(name, idx)
		idx++
	}

	gl.ActiveTexture(0)
}

func (r *Rasterizator) after() {

}


var rect *Mesh
func (r *Rasterizator) DrawRect() {
	r.before()

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

	rect.Draw()

	r.after()
}

// DrawMesh ...
func (r *Rasterizator) Draw(drawable Drawable) {
	r.before()

	time := float32(getTime() - getStartTime()) / 1000
	r.shader.Set1f("iTime", time)
	drawable.Draw()

	r.after()
}



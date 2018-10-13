package render

import (
	// "fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
	id uint32
}

func New2DTexture(width uint, height uint) *Texture {
	t := Texture{}

	gl.GenTextures(1, &t.id)

	gl.BindTexture(gl.TEXTURE_2D, t.id)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, int32(width), int32(height), 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	checkOpenGLError()

	return &t
}

func (t *Texture) GetId() uint32 {
	return t.id
}


func (t *Texture) Destroy() {
	gl.DeleteTextures(1, &t.id)
}

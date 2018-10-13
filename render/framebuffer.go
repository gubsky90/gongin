package render

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Framebuffer struct {
	fbo uint32
	Color *Texture
}

func NewFramebuffer(width uint, height uint) *Framebuffer {
	fb := Framebuffer{}

	gl.GenFramebuffers(1, &fb.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo)

	fb.Color = New2DTexture(width, height)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, fb.Color.GetId(), 0)

	if s := gl.CheckFramebufferStatus(gl.FRAMEBUFFER); s != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Errorf("Bad framebuffer: code %d", s))
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	checkOpenGLError()

	return &fb
}

func (fb *Framebuffer) Distroy() {
	gl.DeleteFramebuffers(1, &fb.fbo)
}

func (fb *Framebuffer) SetAsCurrentRenderTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo)
}
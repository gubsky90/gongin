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

	gl.BindFramebuffer(gl.FRAMEBUFFER, currentFramebuffer)

	checkOpenGLError()

	return &fb
}

func (fb *Framebuffer) Distroy() {
	gl.DeleteFramebuffers(1, &fb.fbo)
}

func (fb *Framebuffer) SetAsCurrentRenderTarget() {
	if currentFramebuffer != fb.fbo {
		currentFramebuffer = fb.fbo
		gl.BindFramebuffer(gl.FRAMEBUFFER, fb.fbo)
	}
}

func (fb *Framebuffer) Clear() {
	fb.SetAsCurrentRenderTarget()

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

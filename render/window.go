package render

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Window struct {
	*glfw.Window
}

func NewWindow (width uint, height uint) *Window {
	var err error
	window := Window{}

	window.Window, err = glfw.CreateWindow(int(width), int(height), "gongin", nil, nil)
	if err != nil {
		panic(fmt.Errorf("Could not create window: %v", err))
	}
	return &window
}

func (w *Window) SetAsCurrentRenderTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

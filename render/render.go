package render

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Render struct {
	win *glfw.Window
}

func New() *Render {
	r := Render{}
	r.win = initGLFW()
	return &r
}

func (r *Render) Destroy() {

}

func (r *Render) ShouldClose() bool {
	return r.win.ShouldClose()
}

func (r *Render) On(name string, cb func()) {

}

func (r *Render) Clear() {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *Render) SwapBuffers() {
	r.win.SwapBuffers()
	glfw.PollEvents()
}

func initGLFW() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("Could not initialize glfw: %v", err))
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

	window, err := glfw.CreateWindow(640, 480, "gongin", nil, nil)
	if err != nil {
		panic(fmt.Errorf("Could not create window: %v", err))
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("Could not initialize OpenGL: %v", err))
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	return window
}

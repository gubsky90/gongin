package render

import (
	"fmt"
	"time"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Render struct {
	win *glfw.Window
	shader *Shader
}

func New() *Render {
	r := Render{}
	r.win = initGLFW()
	return &r
}

func (r *Render) Close() {

}

func (r *Render) ShouldClose() bool {
	return r.win.ShouldClose()
}

func (r *Render) On(name string, cb func()) {

}

func (r *Render) UseShader(shader *Shader) {
	r.shader = shader
}

func getTime() float32 {
	r := time.Now().UnixNano() / int64(time.Second)

	(int32[2])(r)

	return float32(float64(r) / 1000000)
}


func (r *Render) Draw(mesh *Mesh) {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	fmt.Printf("%f\n", getTime())

	if r.shader != nil {
		r.shader.Use()
		// r.shader.Set1f("iTime", float32(time))
		// r.shader.Set1f("iTime", time)
		r.shader.Set3f("iColor", 1.0, 0.2, 0.1)
	}

	mesh.Draw()

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
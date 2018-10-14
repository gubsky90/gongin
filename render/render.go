package render

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type RenderTarget interface {
	SetAsCurrentRenderTarget()
	Clear()
}

type Config struct {
	Width uint
	Height uint
}

type Render struct {
	config Config
	window *Window
}

func New(config Config) *Render {
	
	r := Render{}
	r.config = config
	
	initGLFW()
	r.window = NewWindow(config.Width, config.Height)
	r.window.MakeContextCurrent()
	initOpenGL()
	return &r
}

func (r *Render) Destroy() {

}

var cMain chan func()
func callFromMain(fn func()) {
	if cMain == nil {
		cMain =make(chan func())
	}
	cMain <- fn
}

func doMainFunctions() {
	select {
	case fn := <- cMain:
		fn()
	default:
	}
}

func (r *Render) SwapBuffers() {
	r.window.SwapBuffers()
	glfw.PollEvents()
	doMainFunctions()
}

func (r *Render) GetWindow() *Window {
	return r.window
}

func (r *Render) ShouldClose() bool {
	return r.window.ShouldClose()
}

func (r *Render) On(name string, cb func()) {

}

func initGLFW() {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("Could not initialize glfw: %v", err))
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(fmt.Errorf("Could not initialize OpenGL: %v", err))
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
}

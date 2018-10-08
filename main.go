package main

import (
	"runtime"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var points = []float32 {
	-0.5,  0.5,  0.0,
	-0.5, -0.5,  0.0,
	 0.5,  0.5,  0.0,
	 
	 0.5,  0.5,  0.0,
	-0.5, -0.5,  0.0,
	 0.5, -0.5,  0.0,
}

func init(){
	runtime.LockOSThread();
}

func main(){
	win := initGLFW()
	defer glfw.Terminate()

	mesh := NewMesh(points)

	shader := NewShader(ShaderSource{
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
				vec3 col = iColor;
				col *= abs(sin(iTime * 0.001));
				frag_color = vec4(col, 1.0);
			}
		`,
	})

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	var f float32 = 0.0

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		f += 1.0

		shader.Use()
		shader.Set3f("iColor", 1.0, 0.2, 0.1)
		shader.Set1f("iTime", f)

		mesh.Draw()

		win.SwapBuffers()
		glfw.PollEvents()
	}
}
package render

import (
	"fmt"
	"strings"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var currentShader *Shader = nil

type Shader struct {
	id uint32
}

type ShaderSource struct {
	Vertex string
	Fragment string
}

func NewShader(src ShaderSource) *Shader {
	s := Shader{}
	s.id = gl.CreateProgram()

	if src.Vertex != "" {
		vert, err := compileShader(src.Vertex, gl.VERTEX_SHADER)
		if err != nil {
			panic(err)
		}
		gl.AttachShader(s.id, vert)
	}

	if src.Fragment != "" {
		frag, err := compileShader(src.Fragment, gl.FRAGMENT_SHADER)
		if err != nil {
			panic(err)
		}
		gl.AttachShader(s.id, frag)
	}

	gl.LinkProgram(s.id)

	return &s
}

func (s *Shader) Use() {
	if s != currentShader {
		if currentShader != nil {
			currentShader.unbind()
		}
		currentShader = s
		currentShader.bind()
	}
}

func (s *Shader) bind() {
	gl.UseProgram(s.id)
}

func (s *Shader) unbind() {

}

func (s *Shader) Set1f(name string, f1 float32) {
	cname := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.id, cname)
	gl.Uniform1f(location, f1)
}

func (s *Shader) Set1i(name string, f1 int32) {
	cname := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.id, cname)
	gl.Uniform1i(location, f1)
}


func (s *Shader) Set3f(name string, f1 float32, f2 float32, f3 float32) {
	cname := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.id, cname)
	gl.Uniform3f(location, f1, f2, f3)
}

func compileShader(src string, shaderType uint32) (uint32, error) {
	s := gl.CreateShader(shaderType)
	cstr, free := gl.Strs(src + "\x00")
	gl.ShaderSource(s, 1, cstr, nil)
	free()
	gl.CompileShader(s)

	var status int32
	gl.GetShaderiv(s, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(s, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen + 1))
		gl.GetShaderInfoLog(s, logLen, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to compile: %v", log)
	}

	return s, nil
}

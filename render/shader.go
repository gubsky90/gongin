package render

import (
	"fmt"
	"strings"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var currentShader *Shader = nil

type ShaderTextureItem struct {
	pos int32
	glSlot uint32
	texture *Texture
}

type Shader struct {
	id uint32
	uniforms map[int32]interface{}
	textures map[string]ShaderTextureItem
}

type ShaderSource struct {
	Vertex string
	Fragment string
}

func NewShader(src ShaderSource) *Shader {
	s := Shader{}
	s.uniforms = make(map[int32]interface{})
	s.textures = make(map[string]ShaderTextureItem)

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

func NewShaderFromFile(file string) *Shader {
	s := Shader{}
	return &s
}

func NewShaderWatchFile(file string) *Shader {
	s := Shader{}
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

	for location, value := range s.uniforms {
		switch t:= value.(type) {
		case *int32:
			gl.Uniform1i(location, *value.(*int32))
		case *float32:
			gl.Uniform1f(location, *value.(*float32))
		default:
			panic(fmt.Errorf("type unsupport: %T", t))
		}
	}

	for _, item := range s.textures {
		gl.ActiveTexture(item.glSlot)
		gl.BindTexture(gl.TEXTURE_2D, item.texture.GetId())
	}

	gl.ActiveTexture(0)
}

func (s *Shader) unbind() {

}

func (s *Shader) Bind(name string, value interface{}) {
	cname := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.id, cname)
	if location < 0 {
		panic(fmt.Errorf("Not found unifform with name \"%s\"", name))
	}
	switch t:= value.(type) {
	case *Texture:
		gl.Uniform1i(location, s.setTexture(name, value.(*Texture)))
	case *int32:
	case *float32:
		s.uniforms[location] = value
	default:
		panic(fmt.Errorf("type unsupport: %T", t))
	}
}

func (s *Shader) setTexture(name string, texture *Texture) int32 {
	glTextureSlots := []uint32{
		gl.TEXTURE0,
		gl.TEXTURE1,
		gl.TEXTURE2,
	}

	if _, ok := s.textures[name]; !ok {
		pos := len(s.textures)

		if pos >= len(glTextureSlots) {
			panic(fmt.Errorf("All texture slots used"))
		}

		s.textures[name] = ShaderTextureItem{
			texture: texture,
			pos: int32(pos),
			glSlot: glTextureSlots[pos],
		}
	}

	return s.textures[name].pos
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

package render

import (
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/gubsky90/gongin/tools"
)

var currentShader *Shader = nil

type ShaderTextureItem struct {
	pos int32
	glSlot uint32
	texture *Texture
}

type Shader struct {
	id uint32
	uniformLocations map[string]int32
	textures map[string]ShaderTextureItem
	watcher *tools.FileWatcher
}

type ShaderSource struct {
	Vertex string `yaml:"vertex"`
	Fragment string `yaml:"fragment"`
}

func _newShader() *Shader {
	s := Shader{}
	s.uniformLocations = make(map[string]int32)
	s.textures = make(map[string]ShaderTextureItem)
	return &s
}

func NewShader(src ShaderSource) *Shader {
	s := _newShader()
	s.compile(src)
	return s
}

func NewShaderFromFile(file string) (*Shader, error){
	s := _newShader()
	if err := s.load(file); err != nil {
		return nil, err
	}
	return s, nil
}

func NewShaderWatchFile(file string) (s *Shader, err error) {
	s = _newShader()

	if err := s.load(file); err != nil {
		fmt.Errorf("Shader error: %s", err)
	}

	if s.watcher, err = tools.NewFileWatcher(file, func() {
		fmt.Println("Shader file changes", file)
		callFromMain(func(){
			if err := s.load(file); err != nil {
				fmt.Errorf("Shader error: %s", err)
			} else {
				fmt.Println("Shader successfull recompile")
			}
		})
	}); err != nil {
		panic(fmt.Errorf("Error: %v", err))
		return
	}

	return
}

func (s *Shader) load(file string) error {
	var (
		err error
		data []byte
	)

	fmt.Println("Shader before load file", file)

	data, err = ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	fmt.Println("Shader after load file", file)

	shaderSource := ShaderSource{}
	err = yaml.Unmarshal(data, &shaderSource)
	if err != nil {
		return err
	}

	fmt.Println("Shader after parse file", file)

	return s.compile(shaderSource)
}

func (s *Shader) compile(src ShaderSource) error {
	if s.id > 0 {
		gl.DeleteProgram(s.id)
	}

	s.id = gl.CreateProgram()

	if src.Vertex != "" {
		vert, err := compileShader(src.Vertex, gl.VERTEX_SHADER)
		defer gl.DeleteShader(vert)
		if err != nil {
			return err
		}
		gl.AttachShader(s.id, vert)
	}

	if src.Fragment != "" {
		frag, err := compileShader(src.Fragment, gl.FRAGMENT_SHADER)
		defer gl.DeleteShader(frag)
		if err != nil {
			return err
		}
		gl.AttachShader(s.id, frag)
	}

	gl.LinkProgram(s.id)

	return nil
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

	for _, item := range s.textures {
		gl.ActiveTexture(item.glSlot)
		gl.BindTexture(gl.TEXTURE_2D, item.texture.GetId())
	}

	gl.ActiveTexture(0)
}

func (s *Shader) unbind() {

}

func (s *Shader) getUniformLocation(name string) int32 {
	if _, ok := s.uniformLocations[name]; !ok {
		cname := gl.Str(name + "\x00")
		location := gl.GetUniformLocation(s.id, cname)
		if location < 0 {
			panic(fmt.Errorf("Not found uniform with name \"%s\"", name))
		}
		s.uniformLocations[name] = location
	}

	return s.uniformLocations[name]
}

func (s *Shader) Set(name string, value interface{}) {
	location := s.getUniformLocation(name)
	switch v:= value.(type) {
	case *Texture:
		gl.Uniform1i(location, s.setTexture(name, v))
	case int32:
		gl.Uniform1i(location, v)
	case float32:
		gl.Uniform1f(location, v)
	default:
		panic(fmt.Errorf("type unsupport: %T", v))
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

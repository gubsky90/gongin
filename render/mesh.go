package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	// "github.com/tbogdala/assimp-go"
)

type Mesh struct {
	vbo uint32
	vao uint32
	count int32
}

func NewMesh (points []float32) *Mesh {
	m := Mesh{}

	m.count = int32(len(points) / 3)

	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4 * len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return &m
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.count)
}

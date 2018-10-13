package render

import (
	// "os"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	assimp "github.com/tbogdala/assimp-go"
)

type Mesh struct {
	vbo uint32
	vao uint32
	ebo uint32
	count int32
}

func NewMeshFromFile(file string) *Mesh {
	meshes, err := assimp.ParseFile(file)
	if err != nil {
		panic(fmt.Errorf("Could load file %s: %v", file, err))
	}

	if len(meshes) > 1 {
		panic(fmt.Errorf("File containt more than one mesh. Only one support!"))
	}

	mesh := meshes[0]

	vertices := make([]float32, mesh.VertexCount * 3)
	for i := uint32(0); i < mesh.VertexCount; i++ {
		copy(vertices[i*3:i*3+3], mesh.Vertices[i][0:3])
	}

	indices := make([]uint32, mesh.FaceCount*3)
	for i := uint32(0); i < mesh.FaceCount; i++ {
		copy(indices[i*3:i*3+3], mesh.Faces[i][0:3])
	}

	return NewMesh(vertices, indices)
}

func NewMesh (vertices []float32, indices []uint32) *Mesh {
	m := Mesh{}

	m.count = int32(len(indices))

	gl.GenVertexArrays(1, &m.vao)
	gl.GenBuffers(1, &m.vbo)
	gl.GenBuffers(1, &m.ebo)
	
	gl.BindVertexArray(m.vao)
		// Vertices
		gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

		// Indices
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

		// Attributes
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
		gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)
	return &m
}

func (m *Mesh) Draw() {
	// gl.BindVertexArray(m.vao)
	// gl.DrawArrays(gl.TRIANGLES, 0, m.count)

	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, m.count, gl.UNSIGNED_INT, nil)
}

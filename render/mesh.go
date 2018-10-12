package render

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	assimp "github.com/tbogdala/assimp-go"
)

type Mesh struct {
	vbo uint32
	vao uint32
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

	raw := make([]float32, mesh.FaceCount * 3 * 3)
	for i := uint32(0); i < mesh.FaceCount; i++ {
		offset := i*9
		face := mesh.Faces[i]

		// copy(raw[offset+0:offset+2], mesh.Vertices[face[0]][0:2])
		// copy(raw[offset+3:offset+5], mesh.Vertices[face[1]][0:2])
		// copy(raw[offset+6:offset+8], mesh.Vertices[face[2]][0:2])

		raw[offset + 0] = mesh.Vertices[face[0]][0]
		raw[offset + 1] = mesh.Vertices[face[0]][1]
		raw[offset + 2] = mesh.Vertices[face[0]][2]

		raw[offset + 3] = mesh.Vertices[face[1]][0]
		raw[offset + 4] = mesh.Vertices[face[1]][1]
		raw[offset + 5] = mesh.Vertices[face[1]][2]

		raw[offset + 6] = mesh.Vertices[face[2]][0]
		raw[offset + 7] = mesh.Vertices[face[2]][1]
		raw[offset + 8] = mesh.Vertices[face[2]][2]
	}

	return NewMesh(raw)
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

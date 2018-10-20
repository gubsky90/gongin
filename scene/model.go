package scene

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gubsky90/gongin/render"
)

type Model struct {
	file string

	mesh *render.Mesh
	rotMat mgl32.Mat4
	scaleMat mgl32.Mat4
	transMat mgl32.Mat4
	offsetMat mgl32.Mat4
	sumMat mgl32.Mat4
}

func NewModel(file string) (*Model) {
	m := &Model{
		file: file,
	}
	m.Load()
	return m
}

func (m *Model) calcSumMat() {
	m.sumMat = m.transMat.Mul4(m.rotMat).Mul4(m.offsetMat).Mul4(m.scaleMat)
}

func (m *Model) Load() {
	m.mesh = render.NewMeshFromFile(m.file)
}

func (m *Model) Draw() {
	m.mesh.Draw()
}

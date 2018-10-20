package scene

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gubsky90/gongin/render"
)

type Camera struct {
	shader *render.Shader

	position mgl32.Vec3
	rotation mgl32.Vec3

	projectionMat mgl32.Mat4
	positionMat mgl32.Mat4
	rotaionMat mgl32.Mat4
	resultMat mgl32.Mat4
}

func NewCamera(shader *render.Shader) (camera *Camera) {
	camera = &Camera{
		shader: shader,
	}

	camera.projectionMat = mgl32.Ident4()
	camera.positionMat = mgl32.Ident4()
	camera.rotaionMat = mgl32.Ident4()

	return
}

type Perspective struct {
	Fovy float32
	Aspect float32 // width / height
	Near float32
	Far float32
}

type Ortho struct {
	Left float32
	Right float32
	Bottom float32
	Top float32
	Near float32
	Far float32
}

func (camera *Camera) SetProjection(projection interface{}) {
	switch conf := projection.(type) {
	case Perspective:
		camera.projectionMat = mgl32.Perspective(
			mgl32.DegToRad(conf.Fovy),
			conf.Aspect,
			conf.Near,
			conf.Far,
		)
	case Ortho:
		camera.projectionMat = mgl32.Ortho(
			conf.Left,
			conf.Right,
			conf.Bottom,
			conf.Top,
			conf.Near,
			conf.Far,
		)
	default:
		panic("Bad projection type")	
	}
	camera.updateResultMat()
}

func (camera *Camera) Move(vec mgl32.Vec3) {
	camera.position = camera.position.Add(mgl32.AnglesToQuat(camera.rotation[1], camera.rotation[0], camera.rotation[2], mgl32.XYX).Inverse().Rotate(vec))
	camera.positionMat = mgl32.Translate3D(camera.position[0], camera.position[1], camera.position[2])
	camera.updateResultMat()
}

func (camera *Camera) Rotate(x, y, z float32) {
	camera.rotation[0] += x * 0.01
	camera.rotation[1] += y * 0.01

	camera.rotaionMat = mgl32.AnglesToQuat(camera.rotation[1], camera.rotation[0], camera.rotation[2], mgl32.XYX).Mat4()

	camera.updateResultMat()
}

type HasPosition interface {
	GetPosition() mgl32.Vec3
}

type HasChildren interface {
	GetChildren() []interface{}
}

type IsDrawable interface {
	Draw()
}

func (camera *Camera) Draw(node interface{}) {
	camera.shader.Use()
	camera.draw(node)
}

func (camera *Camera) draw(node interface{}) {
	if node, ok := node.(HasPosition); ok {
		if !camera.CheckSee(node) {
			return
		}
	}

	if node, ok := node.(IsDrawable); ok {
		camera.shader.Set("iCamera", camera.resultMat)
		camera.shader.Set("iModel", mgl32.Ident4())
		node.Draw()
	}

	if node, ok := node.(HasChildren); ok {
		for _, child := range node.GetChildren() {
			camera.Draw(child)
		}
	}
}

func (camera *Camera) updateResultMat() {
	// QuatLookAtV
	camera.resultMat = camera.projectionMat.Mul4(camera.rotaionMat).Mul4(camera.positionMat)
}

func (camera *Camera) CheckSee(node HasPosition) bool {
	return true
}
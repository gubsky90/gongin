package gongin

import (
	"time"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/gubsky90/gongin/render"
	"github.com/gubsky90/gongin/scene"
)

type Config struct {

}

type Gongin struct {
	handlers map[string][]func()
}

func New(conf Config) *Gongin {
	g := Gongin{}
	g.handlers = make(map[string][]func())
	return &g
}

func (g *Gongin) On(name string, cb func()) {
	if _, ok := g.handlers[name]; !ok {
		g.handlers[name] = make([]func(), 0)
	}
	g.handlers[name] = append(g.handlers[name], cb)
}

func (g *Gongin) fire(name string) {
	if cbs, ok := g.handlers[name]; ok {
		for _, cb := range cbs {
			cb()
		}
	}
}

func (g *Gongin) Run2() {
	width := 640
	height := 480

	r := render.New(render.Config{
		Width: width,
		Height: height,
	})
	defer r.Destroy()

	meshShader, err := render.NewShaderFromFile("../assets/main.yml")
	if err != nil {
		panic(err)
	}

	root := scene.NewRoot()
	camera := scene.NewCamera(meshShader)
	camera.SetProjection(scene.Perspective{
		Fovy: 45.0,
		Aspect: float32(width) / float32(height),
		Near: 0.1,
		Far: 1000.0,
	})

	root.Append(scene.NewModel("../assets/teapot.obj"))
	root.Append(scene.NewModel("../assets/box.obj"))

	fmt.Printf("%#v\n", root)

	window := r.GetWindow()

	window.SetAsCurrentRenderTarget()
	window.SetCursorPos(0, 0)

	ticker := time.NewTicker(time.Second / 60)
	for !r.ShouldClose() {
		select {
		case <- ticker.C:
			x, y := window.GetCursorPos()
			window.SetCursorPos(0, 0)

			camera.Rotate(float32(x), float32(y), 0.0)	

			if window.GetKey(glfw.KeyS) == glfw.Press {
				camera.Move(mgl32.Vec3{0.0, 0.0, -0.1})
			}
			if window.GetKey(glfw.KeyW) == glfw.Press {
				camera.Move(mgl32.Vec3{0.0, 0.0, 0.1})
			}
			if window.GetKey(glfw.KeyD) == glfw.Press {
				camera.Move(mgl32.Vec3{-0.1, 0.0, 0.0})
			}
			if window.GetKey(glfw.KeyA) == glfw.Press {
				camera.Move(mgl32.Vec3{0.1, 0.0, 0.0})
			}

			if window.GetKey(glfw.KeyEscape) == glfw.Press {
				return
			}

			window.Clear()
			camera.Draw(root)
			r.SwapBuffers()
		}
	}
}

func (g *Gongin) Run() {
	width := 640
	height := 480

	r := render.New(render.Config{
		Width: width,
		Height: height,
	})
	defer r.Destroy()

	fb := render.NewFramebuffer(width, height)
	win := r.GetWindow()

	meshShader, err := render.NewShaderFromFile("../assets/main.yml")
	if err != nil {
		panic(err)
	}

	postShader, err := render.NewShaderFromFile("../assets/post.yml")
	if err != nil {
		panic(err)
	}

	meshRaster := render.NewRasterizator(fb, meshShader)
	postRaster := render.NewRasterizator(win, postShader)
	postRaster.Set("screenTexture", fb.Color)

	mesh := render.NewMeshFromFile("../assets/teapot.obj")
	g.fire("ready")

	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)

	iProjection := mgl32.Perspective(
		mgl32.DegToRad(45.0),
		float32(width)/float32(height),
		0.1,
		10.0,
	)

	iCamera := mgl32.LookAtV(
		mgl32.Vec3{0, 0, -3},
		mgl32.Vec3{0, 0, 0},
		mgl32.Vec3{0, 1, 0},
	)

	for !r.ShouldClose() {

		// start := time.Now()

		iTime := float32(getTime() - getStartTime()) / 1000

		angle := iTime * 1.0

		modelTrans := mgl32.Translate3D(0.5, 0.0, 0.0)
		modelAxisTrans := mgl32.Translate3D(0.0, 0.0, 0.0)
		modelRot := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0.0, 1.0, 0.0})
		modelScale := mgl32.Scale3D(0.1, 0.1, 0.1)

		iModel := modelTrans.Mul4(modelRot).Mul4(modelAxisTrans).Mul4(modelScale)

		fb.Clear()
		meshRaster.Set("iTime", iTime)
		meshRaster.Set("iModel", iModel)
		meshRaster.Set("iCamera", iCamera)
		meshRaster.Set("iProjection", iProjection)
		meshRaster.Draw(mesh)

		// win.Clear()
		postRaster.Set("iTime", iTime)
		postRaster.DrawRect()

		r.SwapBuffers()
		// fmt.Printf("Frame time: %s\n", time.Since(start))
	}
}

var startTime int64

func getStartTime() int64 {
	if startTime == 0 {
		startTime = getTime()
	}
	return startTime
}

func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
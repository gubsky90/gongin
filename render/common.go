package render

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var currentFramebuffer uint32 = 0

func checkOpenGLError() {
	var errors []uint32
	for {
		err := gl.GetError()
		if err == gl.NO_ERROR {
			break;
		}
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		panic(fmt.Errorf("OpenGL Error: %v", errors))
	}
}

func getMaxDrawBuffers() (maxDrawBuffers int32) {
	gl.GetIntegerv(gl.MAX_DRAW_BUFFERS, &maxDrawBuffers)
	return
}
package render

import (
	"fmt"
	"time"
	"github.com/go-gl/gl/v4.1-core/gl"
)

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
package render

import (
	"time"
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

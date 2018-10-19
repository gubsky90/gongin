package tools

import (
	"os"
	"time"
)

type FileWatcher struct {
	ticker chan struct{}
}

func NewFileWatcher(file string, cb func()) (fw *FileWatcher, err error) {
	var lastModify time.Time

	interval := 500 * time.Millisecond
	debounce := time.Second

	if lastModify, err = getModTime(file); err != nil {
		return
	}

	fw = &FileWatcher{}
	ticker := time.NewTicker(interval)
	fw.ticker = make(chan struct{})
	go func(){
		timer := time.NewTimer(0)
		<- timer.C
		for {
			select {
			case <- ticker.C:
				t, err := getModTime(file)
				if err != nil {
					continue
				}
				if !t.Equal(lastModify) {
					lastModify = t
					timer.Reset(debounce)
				}
			case <- timer.C:
				cb()
			case <- fw.ticker:
				ticker.Stop()
				return
			}
		}
	}()

	return
}

func (fw *FileWatcher) Destroy() {
	close(fw.ticker)
}

func getModTime(file string) (t time.Time, err error) {
	var info os.FileInfo
	if info, err = os.Stat(file); err == nil {
		t = info.ModTime()
	}
	return 
}

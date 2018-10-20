package main

import (
	"fmt"
	"runtime"
	"github.com/gubsky90/gongin"
)

func main(){
	runtime.LockOSThread();

	fmt.Println("Ok")

	gongin := gongin.New(gongin.Config{

	})

	gongin.On("ready", func(){
		fmt.Println("In ready handler")
	})

	// gongin.Run()
	gongin.Run2()
}
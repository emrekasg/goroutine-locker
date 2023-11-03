package main

import (
	"fmt"

	"time"

	"github.com/emrekasg/goroutine-locker/cpu"
)

func main() {
	coreManager := cpu.NewCoreManager()

	startTime := time.Now()
	defer func() {
		fmt.Println("Total duration:", time.Since(startTime))
	}()

	// run 250.000 task on each core except 4th core
	for i := 0; i < 3; i++ {
		for j := 0; j < 250e3; j++ {
			coreManager.RunTask(i, PrintHelloWorld(j))
		}
	}

	fmt.Println("All tasks are completed")
	generalInfo := coreManager.GetCoreInfo()
	fmt.Printf("%+v", generalInfo)

	// get total goroutine count
	fmt.Println("Total goroutine count:", coreManager.GetGoRoutineCount())
}

func PrintHelloWorld(x int) func() {
	return func() {
		// fmt.Printf("Hello World: %d\n", x)
		return
	}
}

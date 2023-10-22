package main

import (
	"fmt"
	"go-board/logic"
	"sync"
	"time"
)

func main() {
	startTime := time.Now()

	var wg sync.WaitGroup

	wg.Add(1)
	go executeGaltonBoard("./data/config.json", &wg)

	wg.Wait()

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	fmt.Printf("Total time: %v", elapsed)
}

func executeGaltonBoard(route string, wg *sync.WaitGroup) {
	defer wg.Done()

	gb := logic.NewGaltonBoard(route)
	gb.BuildObstacles()
	gb.BuildSpheres()
	//gb.BuildBorders()
	gb.BuildMesh()

	gb.Run()
}

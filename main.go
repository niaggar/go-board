package main

import (
	"fmt"
	"go-board/logic"
	"time"
)

func main() {
	startTime := time.Now()

	gb := logic.NewGaltonBoard("./data/config.json")
	gb.BuildObstacles()
	gb.BuildSpheres()
	gb.BuildBorders()

	gb.Run()

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	fmt.Printf("Total time: %v", elapsed)
}

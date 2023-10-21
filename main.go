package main

import (
	"go-board/logic"
)

func main() {
	gb := logic.NewGaltonBoard("./data/config.json")
	gb.BuildObstacles()
	gb.BuildSpheres()

	gb.Run()
}

package main

import (
	"fmt"
	"go-board/gmath"
	"go-board/logic"
	"go-board/models"
)

func main() {
	fmt.Print("Galton board in Go")

	// Create engine
	gravity := gmath.NewVector(0, -9.8)
	bounds := gmath.NewVector(10, 10)
	damping := 0.9
	timeStep := 1 / 60.0
	subSteps := 8

	engine := logic.NewEngine(gravity, bounds, damping, timeStep, subSteps)

	// Create spheres
	radius := 1.0
	mass := 1.0
	for i := 0; i < 10; i++ {
		x := 5.0
		y := 8.0

		sphere := models.NewSphere(x, y, radius, mass, damping, models.DYNAMIC)
		engine.AddSphere(sphere)
	}

	// Run simulation
	engine.CreateFile()
	for i := 0; i < 1000; i++ {
		engine.Update()
		engine.Export()
	}
	engine.CloseFile()

	fmt.Print("Simulation finished")
}

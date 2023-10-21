package logic

import (
	"go-board/export"
	"go-board/gmath"
	"go-board/logic/config"
	"go-board/logic/physics"
	"go-board/models"
	"math/rand"
)

type GaltonBoard struct {
	engine        physics.Engine
	exporter      *export.Exporter
	newConfig     *config.NewConfig
	currentConfig *config.CurrentConfig
	maxTime       float64
	currentTime   float64
	timeStep      float64
	borders       []models.Sphere
}

func NewGaltonBoard(route string) *GaltonBoard {
	configuration := config.GetNewConfiguration(route)

	gridSize := configuration.Board.GridSize
	xSize := gridSize * float64(configuration.Board.ColumnNumber)
	ySize := gridSize * float64(configuration.Board.RowNumber)

	bounds := gmath.NewVector(xSize, ySize)
	gravity := configuration.Board.Gravity
	damping := configuration.Board.Damping

	maxTime := configuration.Board.MaxTime.MaxTime
	currentTime := 0.0
	timeStep := configuration.Board.TimeStep
	subSteps := configuration.Board.SubSteps

	exporter := export.NewExporter(configuration.Board.ExportPathName, configuration.Board.ExportHistogramName)
	engine := physics.NewEngine(gravity, bounds, damping, timeStep, subSteps, exporter)

	return &GaltonBoard{
		engine:        engine,
		exporter:      exporter,
		newConfig:     &configuration,
		currentConfig: nil,
		maxTime:       maxTime,
		currentTime:   currentTime,
		timeStep:      timeStep,
	}
}

func (gb *GaltonBoard) Run() {
	gb.exporter.CreatePathFile()

	for gb.currentTime < gb.maxTime {
		gb.engine.Update()
		gb.engine.ExportCurrentState()

		gb.currentTime += gb.timeStep
		// gb.exporter.ExportCurrentState(gb.engine.Objects, gb.currentTime)
	}

	gb.exporter.ClosePathFile()
}

func (gb *GaltonBoard) BuildSpheres() {
	if gb.currentConfig != nil {

	} else if gb.newConfig != nil {
		autoCreation := gb.newConfig.CreateBalls.Creation.Enabled
		if autoCreation {
			ballsPoints := gb.buildSpheres()
			for _, ball := range *ballsPoints {
				gb.engine.AddSphere(ball)
			}
		} else {
			ballsPoints := gb.newConfig.CreateBalls.Positions
			for _, ball := range ballsPoints {
				sphere := models.NewSphere(ball.Position.X, ball.Position.Y, ball.Radius, ball.Mass, ball.Damping, models.DYNAMIC)
				sphere.CanCollide = ball.Collision
				gb.engine.AddSphere(sphere)
			}
		}
	}
}

func (gb *GaltonBoard) BuildObstacles() {
	if gb.currentConfig != nil {

	} else if gb.newConfig != nil {
		autoCreation := gb.newConfig.CreateObstacles.Creation.Enabled
		if autoCreation {
			pegsPoints := gb.buildObstaclesCruz()
			for _, peg := range *pegsPoints {
				gb.engine.AddObstacle(peg)
			}
		} else {
			pegsPoints := gb.newConfig.CreateObstacles.Positions
			for _, peg := range pegsPoints {
				sphere := models.NewSphere(peg.Position.X, peg.Position.Y, peg.Radius, peg.Mass, peg.Damping, models.STATIC)
				gb.engine.AddObstacle(sphere)
			}
		}
	}
}

func (gb *GaltonBoard) buildObstaclesCruz() *[]models.Sphere {
	pegsPoints := make([]models.Sphere, 0)

	minSize := gb.newConfig.CreateObstacles.Creation.Radius.Min
	maxSize := gb.newConfig.CreateObstacles.Creation.Radius.Max
	direction := gb.newConfig.CreateObstacles.Creation.Direction

	cols := gb.newConfig.Board.ColumnNumber
	rows := gb.newConfig.Board.RowNumber
	bounds := gb.engine.Bounds
	radius := minSize + rand.Float64()*(maxSize-minSize)
	yGlobalOffset := gb.newConfig.CreateObstacles.Creation.YOffset
	mass := gb.newConfig.CreateObstacles.Creation.Mass
	damping := gb.newConfig.CreateObstacles.Creation.Damping

	xOffset := bounds.X / float64(cols+1)
	yOffset := (bounds.Y - yGlobalOffset) / float64(rows+1)

	for i := 0; i < rows; i++ {
		if direction == 0 {
			radius = minSize + ((maxSize-minSize)/float64(rows))*float64(i)
		}

		for j := 0; j < cols+2; j++ {
			if direction == 1 {
				radius = minSize + ((maxSize-minSize)/float64(cols))*float64(j)
			}

			if i%2 == 0 {
				x := xOffset * float64(j)
				y := yOffset*float64(i+1) + yGlobalOffset

				sphere := models.NewSphere(x, y, radius, mass, damping, models.STATIC)
				pegsPoints = append(pegsPoints, sphere)
			} else {
				if j >= cols-1 {
					continue
				} else {
					x := xOffset*float64(j+1) + xOffset/2
					y := yOffset*float64(i+1) + yGlobalOffset

					sphere := models.NewSphere(x, y, radius, mass, damping, models.STATIC)
					pegsPoints = append(pegsPoints, sphere)
				}
			}
		}
	}

	return &pegsPoints
}

func (gb *GaltonBoard) buildSpheres() *[]models.Sphere {
	ballsPoints := make([]models.Sphere, 0)

	canCollide := gb.newConfig.CreateBalls.Collisions
	mass := gb.newConfig.CreateBalls.Creation.Mass
	damping := gb.newConfig.CreateBalls.Creation.Damping

	rRange := gb.newConfig.CreateBalls.Creation.Radius
	xRange := gb.newConfig.CreateBalls.Creation.Position.X
	vxRange := gb.newConfig.CreateBalls.Creation.Velocity.X
	vyRange := gb.newConfig.CreateBalls.Creation.Velocity.Y

	nSpheres := gb.newConfig.CreateBalls.Creation.Count
	bounds := gb.engine.Bounds

	for i := 0; i < nSpheres; i++ {
		radius := rRange.Min + rand.Float64()*(rRange.Max-rRange.Min)
		x := xRange.Min + rand.Float64()*(xRange.Max-xRange.Min)
		y := bounds.Y - 2*radius

		vx := vxRange.Min + rand.Float64()*(vxRange.Max-vxRange.Min)
		vy := vyRange.Min + rand.Float64()*(vyRange.Max-vyRange.Min)

		sphere := models.NewSphere(x, y, radius, mass, damping, models.DYNAMIC)
		sphere.CanCollide = canCollide
		sphere.Velocity = gmath.NewVector(vx, vy)

		ballsPoints = append(ballsPoints, sphere)
	}

	return &ballsPoints
}

func (gb *GaltonBoard) buildBorders() {

}

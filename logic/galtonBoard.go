package logic

import (
	"go-board/export"
	"go-board/gmath"
	"go-board/logic/config"
	"go-board/models"
	"math/rand"
)

type GaltonBoard struct {
	engine        *Engine
	exporter      *export.Exporter
	newConfig     *config.NewConfig
	currentConfig *config.CurrentConfig
}

func NewGaltonBoard(route string) *GaltonBoard {
	configuration := config.GetNewConfiguration(route)

	gridSize := configuration.Board.GridSize
	xSize := gridSize * float64(configuration.Board.ColumnNumber)
	ySize := gridSize * float64(configuration.Board.RowNumber)

	bounds := gmath.NewVector(xSize, ySize)
	gravity := configuration.Board.Gravity
	damping := configuration.Board.Damping
	timeStep := configuration.Board.TimeStep
	subSteps := configuration.Board.SubSteps
	exportRoute := configuration.Board.ExportRoute

	exporter := export.NewExporter(exportRoute)
	engine := NewEngine(gravity, bounds, damping, timeStep, subSteps, exporter)

	return &GaltonBoard{
		engine:        engine,
		exporter:      exporter,
		newConfig:     &configuration,
		currentConfig: nil,
	}
}

func (gb *GaltonBoard) Init() {
	// Create engine
	// Create spheres
	// Create exporter
}

func (gb *GaltonBoard) Run() {
	// Run engine
	// Run exporter
}

func (gb *GaltonBoard) Stop() {
	// Stop engine
	// Stop exporter
}

func (gb *GaltonBoard) BuildObstacles() {
	if gb.currentConfig != nil {

	} else if gb.newConfig != nil {
		autoCreation := gb.newConfig.CreateObstacles.Creation.Enabled
		if autoCreation {

		} else {
			// Manual creation
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

	xOffset := bounds.X() / float64(cols+1)
	yOffset := (bounds.Y() - yGlobalOffset) / float64(rows+1)

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
				pegsPoints = append(pegsPoints, *sphere)
			} else {
				if j >= cols-1 {
					continue
				} else {
					x := xOffset*float64(j+1) + xOffset/2
					y := yOffset*float64(i+1) + yGlobalOffset

					sphere := models.NewSphere(x, y, radius, mass, damping, models.STATIC)
					pegsPoints = append(pegsPoints, *sphere)
				}
			}
		}
	}

	return &pegsPoints
}

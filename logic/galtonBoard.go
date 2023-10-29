package logic

import (
	"go-board/logic/physics"
	"go-board/models"
	"go-board/utils/export"
	"go-board/utils/gmath"
)

type GaltonBoard struct {
	*models.SimulationProps
	*models.ExperimentProps
	*models.BoardProps

	Bounds   *models.Bounds
	Engine   *physics.Engine
	Exporter *export.Exporter

	Balls     []*models.Ball
	Obstacles []*models.Ball
	Borders   []*gmath.Vector
}

func NewGaltonBoard(config *models.ConfigurationFile) *GaltonBoard {
	bounds := models.Bounds{
		Width:  float32(config.Board.ColumnNumber) * config.Board.CellSize,
		Height: float32(config.Board.RowNumber) * config.Board.CellSize,
	}

	exporter := export.NewExporter(config.Experiment.ExportPaths, config.Experiment.ExportHistogram)
	engine := physics.NewEngine(&config.Simulation, &config.Board, &bounds)

	return &GaltonBoard{
		Engine:          engine,
		Exporter:        exporter,
		Bounds:          &bounds,
		BoardProps:      &config.Board,
		SimulationProps: &config.Simulation,
		ExperimentProps: &config.Experiment,
		//Borders:         BuildBorders(&bounds),
		Balls:     BuildSpheres(&config.Balls, &bounds),
		Obstacles: BuildObstacles(&config.Obstacles, &config.Board, &bounds),
	}
}

func (gb *GaltonBoard) Run() {
	currentTime := float32(0)
	currentFrame := 0
	exportFrame := 0

	// Add obstacles
	gb.AddObstacle()

	for currentTime < gb.Stop.MaxTime {
		// Add balls
		gb.AddBall(currentFrame)

		// Update engine
		gb.Engine.Update()

		// Export path
		if exportFrame >= gb.FrameRate {
			exportFrame = 0
			gb.ExportPath()
		}

		// Validate stop
		stop := gb.ValidateStop(currentTime)
		if stop {
			break
		}

		currentTime += gb.TimeStep
		currentFrame++
		exportFrame++
	}

	gb.ExportHistogram()
}

func (gb *GaltonBoard) Finish() {
	gb.Exporter.CloseFiles()
}

func (gb *GaltonBoard) AddBall(currentFrame int) {
	currentCount := len(gb.Engine.Balls)
	totalCount := len(gb.Balls)

	if gb.CreationFrame == 0 {
		if currentCount == 0 {
			for i := 0; i < totalCount; i++ {
				ballToAdd := gb.Balls[currentCount]
				gb.Engine.AddBall(ballToAdd)
			}
		}

		return
	}

	addFrame := currentFrame % gb.CreationFrame
	if addFrame == 0 {
		totalToAdd := gb.CreationNum

		if currentCount < totalCount {
			if currentCount+gb.CreationNum > totalCount {
				totalToAdd = totalCount - (currentCount + gb.CreationNum)
			}

			for i := 0; i < totalToAdd; i++ {
				ballToAdd := gb.Balls[currentCount+i]
				gb.Engine.AddBall(ballToAdd)
			}
		}
	}
}

func (gb *GaltonBoard) AddObstacle() {
	currentCount := len(gb.Engine.Obstacles)
	totalCount := len(gb.Obstacles)

	if currentCount < totalCount {
		for i := 0; i < totalCount; i++ {
			obsToAdd := gb.Obstacles[i]
			gb.Engine.AddObstacle(obsToAdd)
		}
	}
}

func (gb *GaltonBoard) ValidateStop(currentTime float32) bool {
	if currentTime >= gb.Stop.MaxTime {
		return true
	}

	stopCount := 0
	for i := 0; i < len(gb.Balls); i++ {
		ball := gb.Balls[i]
		if !ball.Active {
			stopCount++
		}
	}

	if stopCount >= len(gb.Balls) {
		return true
	}

	return false
}

func (gb *GaltonBoard) ExportPath() {
	gb.Exporter.ExportCurrentState(gb.Engine.Balls, gb.Engine.Obstacles, gb.Borders)
}

func (gb *GaltonBoard) ExportHistogram() {
	gb.Exporter.ExportHistogram(gb.Balls, gb.ColumnNumber, gb.CellSize)
}

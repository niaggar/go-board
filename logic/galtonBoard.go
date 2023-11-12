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

	currentTime  float32
	currentFrame int
	exportFrame  int

	Finished bool
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
		Borders:         BuildBorders(&bounds),
		Balls:           BuildSpheres(&config.Balls, &bounds),
		Obstacles:       BuildObstacles(&config.Obstacles, &config.Board, &bounds),
	}
}

func (gb *GaltonBoard) RunAll() {
	gb.currentTime = 0
	gb.currentFrame = 0
	gb.exportFrame = 0

	// Add obstacles
	gb.AddObstacle()

	dt := gb.TimeStep
	for gb.currentTime < gb.Stop.MaxTime && !gb.Finished {
		gb.RunStep(dt)

		if gb.Finished {
			break
		}
	}
}

func (gb *GaltonBoard) RunStep(dt float32) {
	if gb.currentTime < gb.Stop.MaxTime && !gb.Finished {
		// Add balls
		gb.AddBall(gb.currentFrame)

		// Update engine
		gb.Engine.Update(dt)

		// Export path
		if gb.exportFrame >= gb.FrameRate {
			gb.exportFrame = 0
			gb.ExportPath()
		}

		// Validate stop
		stop := gb.ValidateStop(gb.currentTime)
		if stop {
			gb.Finished = true
		}

		gb.currentTime += gb.TimeStep
		gb.currentFrame++
		gb.exportFrame++
	} else {
		gb.Finished = true
	}
}

func (gb *GaltonBoard) Finish() {
	gb.ExportHistogram()

	gb.currentFrame = 0
	gb.currentTime = 0
	gb.exportFrame = 0
	gb.Finished = false

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

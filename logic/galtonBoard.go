package logic

import (
	"go-board/export"
	"go-board/gmath"
	"go-board/logic/config"
	"go-board/logic/physics"
	"go-board/models"
	"math"
	"math/rand"
)

type GaltonBoard struct {
	engine            physics.Engine
	pathExporter      *export.Exporter
	histogramExporter *export.Exporter
	config            *config.NewConfig
	maxTime           float32
	currentTime       float32
	timeStep          float32
	borders           []*gmath.Vector
}

func NewGaltonBoard(route string) *GaltonBoard {
	configuration := config.GetNewConfiguration(route)

	gridSize := configuration.Board.GridSize
	xSize := gridSize * float32(configuration.Board.ColumnNumber)
	ySize := gridSize * float32(configuration.Board.RowNumber)

	bounds := gmath.NewVector(xSize, ySize)
	gravity := configuration.Board.Gravity
	damping := configuration.Board.Damping

	maxTime := configuration.Board.MaxTime.MaxTime
	currentTime := float32(0)
	timeStep := configuration.Board.TimeStep
	subSteps := configuration.Board.SubSteps

	var (
		pathExporter      *export.Exporter = nil
		histogramExporter *export.Exporter = nil
	)

	if configuration.ExportPath.Enabled {
		pathExporter = export.NewExporter(configuration.ExportPath.Path)
	}

	if configuration.ExportHisto.Enabled {
		histogramExporter = export.NewExporter(configuration.ExportHisto.Path)
	}

	engine := physics.NewEngine(gravity, bounds, damping, timeStep, subSteps, pathExporter, histogramExporter)

	return &GaltonBoard{
		engine:            engine,
		pathExporter:      pathExporter,
		histogramExporter: histogramExporter,
		config:            &configuration,
		maxTime:           maxTime,
		currentTime:       currentTime,
		timeStep:          timeStep,
	}
}

func (gb *GaltonBoard) Run() {
	if gb.pathExporter != nil {
		gb.pathExporter.CreateFile()
		defer gb.pathExporter.CloseFile()
	}

	if gb.histogramExporter != nil {
		gb.histogramExporter.CreateFile()
		defer gb.histogramExporter.CloseFile()
	}

	for gb.currentTime < gb.maxTime {
		gb.engine.Update()

		if gb.pathExporter != nil {
			gb.engine.ExportCurrentState(gb.borders)
		}

		gb.currentTime += gb.timeStep
		// gb.pathExporter.ExportCurrentState(gb.engine.Objects, gb.currentTime)
	}
}

func (gb *GaltonBoard) BuildMesh() {
	rObsMax := gb.config.CreateObstacles.Creation.Radius.Max
	rBallMax := gb.config.CreateBalls.Creation.Radius.Max

	globalMax := math.Max(float64(rObsMax), float64(rBallMax))

	gb.engine.MaxSize = float32(globalMax) * 2
	gb.engine.CreateMesh()
}

func (gb *GaltonBoard) BuildSpheres() {
	autoCreation := gb.config.CreateBalls.Creation.Enabled
	if autoCreation {
		ballsPoints := gb.buildSpheres()
		for _, ball := range *ballsPoints {
			gb.engine.AddSphere(ball)
		}
	} else {
		ballsPoints := gb.config.CreateBalls.Positions
		for _, ball := range ballsPoints {
			sphere := models.NewSphere(ball.Position.X, ball.Position.Y, ball.Radius, ball.Mass, ball.Damping, models.DYNAMIC)
			sphere.CanCollide = ball.Collision
			gb.engine.AddSphere(*sphere)
		}
	}
}

func (gb *GaltonBoard) BuildObstacles() {
	autoCreation := gb.config.CreateObstacles.Creation.Enabled
	if autoCreation {
		pegsPoints := gb.buildObstaclesCruz()
		for _, peg := range *pegsPoints {
			gb.engine.AddObstacle(peg)
		}
	} else {
		pegsPoints := gb.config.CreateObstacles.Positions
		for _, peg := range pegsPoints {
			sphere := models.NewSphere(peg.Position.X, peg.Position.Y, peg.Radius, peg.Mass, peg.Damping, models.STATIC)
			gb.engine.AddObstacle(*sphere)
		}
	}
}

func (gb *GaltonBoard) BuildBorders() {
	width := gb.engine.Bounds.X
	height := gb.engine.Bounds.Y

	borders := make([]*gmath.Vector, 0)

	for borderDivision := float32(0); borderDivision < height; borderDivision += (0.5 * 2) {
		v1 := gmath.NewVector(0, borderDivision)
		v2 := gmath.NewVector(width, borderDivision)

		borders = append(borders, &v1, &v2)
	}
	for borderDivision := float32(0); borderDivision < width; borderDivision += (0.5 * 2) {
		v1 := gmath.NewVector(borderDivision, 0)
		v2 := gmath.NewVector(borderDivision, height)

		borders = append(borders, &v1, &v2)
	}

	gb.borders = borders
}

func (gb *GaltonBoard) buildObstaclesCruz() *[]models.Sphere {
	pegsPoints := make([]models.Sphere, 0)

	minSize := gb.config.CreateObstacles.Creation.Radius.Min
	maxSize := gb.config.CreateObstacles.Creation.Radius.Max
	direction := gb.config.CreateObstacles.Creation.Direction

	cols := gb.config.Board.ColumnNumber
	rows := gb.config.Board.RowNumber
	bounds := gb.engine.Bounds
	radius := minSize + rand.Float32()*(maxSize-minSize)
	yGlobalOffset := gb.config.CreateObstacles.Creation.YOffset
	mass := gb.config.CreateObstacles.Creation.Mass
	damping := gb.config.CreateObstacles.Creation.Damping

	xOffset := bounds.X / float32(cols+1)
	yOffset := (bounds.Y - yGlobalOffset) / float32(rows+1)

	for i := 0; i < rows; i++ {
		if direction == 0 {
			radius = minSize + ((maxSize-minSize)/float32(rows))*float32(i)
		}

		for j := 0; j < cols+2; j++ {
			if direction == 1 {
				radius = minSize + ((maxSize-minSize)/float32(cols))*float32(j)
			}

			if i%2 == 0 {
				x := xOffset * float32(j)
				y := yOffset*float32(i+1) + yGlobalOffset

				sphere := models.NewSphere(x, y, radius, mass, damping, models.STATIC)
				pegsPoints = append(pegsPoints, *sphere)
			} else {
				if j >= cols-1 {
					continue
				} else {
					x := xOffset*float32(j+1) + xOffset/2
					y := yOffset*float32(i+1) + yGlobalOffset

					sphere := models.NewSphere(x, y, radius, mass, damping, models.STATIC)
					pegsPoints = append(pegsPoints, *sphere)
				}
			}
		}
	}

	return &pegsPoints
}

func (gb *GaltonBoard) buildSpheres() *[]models.Sphere {
	ballsPoints := make([]models.Sphere, 0)

	canCollide := gb.config.CreateBalls.Collisions
	mass := gb.config.CreateBalls.Creation.Mass
	damping := gb.config.CreateBalls.Creation.Damping

	rRange := gb.config.CreateBalls.Creation.Radius
	xRange := gb.config.CreateBalls.Creation.Position.X
	vxRange := gb.config.CreateBalls.Creation.Velocity.X
	vyRange := gb.config.CreateBalls.Creation.Velocity.Y

	nSpheres := gb.config.CreateBalls.Creation.Count
	bounds := gb.engine.Bounds

	for i := 0; i < nSpheres; i++ {
		radius := rRange.Min + rand.Float32()*(rRange.Max-rRange.Min)
		x := xRange.Min + rand.Float32()*(xRange.Max-xRange.Min)
		y := bounds.Y - 2*radius

		vx := vxRange.Min + rand.Float32()*(vxRange.Max-vxRange.Min)
		vy := vyRange.Min + rand.Float32()*(vyRange.Max-vyRange.Min)

		sphere := models.NewSphere(x, y, radius, mass, damping, models.DYNAMIC)
		sphere.CanCollide = canCollide
		sphere.Velocity = gmath.NewVector(vx, vy)

		ballsPoints = append(ballsPoints, *sphere)
	}

	return &ballsPoints
}

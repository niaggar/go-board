package physics

import (
	"go-board/models"
	"go-board/utils"
	"sync"
)

type Engine struct {
	*models.SimulationProps
	*models.BoardProps
	Balls     []*models.Ball
	Obstacles []*models.Ball
	Bounds    *models.Bounds
	Mesh      *models.Mesh
	Dt        float32
	Ended     bool
}

func NewEngine(simProps *models.SimulationProps, boardPros *models.BoardProps, bounds *models.Bounds) *Engine {
	return &Engine{
		SimulationProps: simProps,
		BoardProps:      boardPros,
		Bounds:          bounds,
		Balls:           make([]*models.Ball, 0),
		Obstacles:       make([]*models.Ball, 0),
		Dt:              simProps.TimeStep / float32(simProps.SubSteps),
		Mesh:            models.NewMesh(*bounds, simProps.MaxSize),
	}
}

func (e *Engine) AddBall(s *models.Ball) {
	e.Balls = append(e.Balls, s)
}

func (e *Engine) AddObstacle(s *models.Ball) {
	e.Obstacles = append(e.Obstacles, s)
}

func (e *Engine) Update(dt float32) {
	for i := 0; i < e.SubSteps; i++ {
		e.updateBodiesParallel(dt)
		e.validateCollisionsParallel()
	}
}

func (e *Engine) updateBodiesParallel(dt float32) {
	var wg sync.WaitGroup

	for i := 0; i < len(e.Obstacles); i++ {
		e.Obstacles[i].Update(dt)
	}

	for i := 0; i < len(e.Balls); i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			e.Balls[i].ApplyForce(&e.Gravity)
			e.Balls[i].Update(dt)
			CollisionBounds(e.Balls[i], e.BoardProps, e.Bounds, e.Stop.TouchGround)
		}(i)
	}

	wg.Wait()
}

func (e *Engine) validateCollisionsParallel() {
	e.updateMesh()

	verDivSize := e.Mesh.Rows / utils.GridParallelDivisionVertical
	horDivSize := e.Mesh.Columns / utils.GridParallelDivisionHorizontal

	var wg sync.WaitGroup

	for sectVer := 0; sectVer <= utils.GridParallelDivisionVertical; sectVer++ {
		for sectHor := 0; sectHor <= utils.GridParallelDivisionHorizontal; sectHor++ {
			startX := sectHor * horDivSize
			startY := sectVer * verDivSize

			wg.Add(1)
			go e.collisionsCellsAround(startX, startY, horDivSize, verDivSize, &wg)
		}
	}

	wg.Wait()
	e.Mesh.Clear()
}

func (e *Engine) collisionsCellsAround(startX, startY, horDivSize, verDivSize int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := startX; i < startX+horDivSize; i++ {
		for j := startY; j < startY+verDivSize; j++ {
			ballsIndex, obstaclesIndex := e.Mesh.GetElementsAround(i, j)

			for k := 0; k < len(ballsIndex); k++ {
				ball := e.Balls[*ballsIndex[k]]
				if !ball.Active {
					continue
				}

				for l := 0; l < len(obstaclesIndex); l++ {
					obstacle := e.Obstacles[*obstaclesIndex[l]]
					if !obstacle.Active {
						continue
					}

					ValidateCollision(ball, obstacle)
				}

				// Validate that collisions between balls are activated
				if !e.Collisions {
					continue
				}

				for m := k + 1; m < len(ballsIndex); m++ {
					otherBall := e.Balls[*ballsIndex[m]]
					if !otherBall.Active {
						continue
					}

					ValidateCollision(ball, otherBall)
				}
			}
		}
	}
}

func (e *Engine) updateMesh() {
	for i := 0; i < len(e.Balls); i++ {
		ball := e.Balls[i]
		if ball.Active {
			e.Mesh.AddBall(ball.Position, i)
		}
	}

	for i := 0; i < len(e.Obstacles); i++ {
		e.Mesh.AddObstacle(e.Obstacles[i].Position, i)
	}
}

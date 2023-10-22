package physics

import (
	"go-board/export"
	"go-board/gmath"
	"go-board/models"
	"sync"
)

type Engine struct {
	Gravity   gmath.Vector
	Objects   []*models.Sphere
	Obstacles []*models.Sphere
	Damping   float32
	TimeStep  float32
	SubSteps  int
	Dt        float32
	Bounds    gmath.Vector
	PathExp   *export.Exporter
	HisExp    *export.Exporter
	Mesh      *models.Mesh
	MaxSize   float32
}

func NewEngine(gravity, bounds gmath.Vector, damping, timeStep float32, subSteps int, pathExp, hisExp *export.Exporter) Engine {
	return Engine{
		Gravity:   gravity,
		Objects:   make([]*models.Sphere, 0),
		Obstacles: make([]*models.Sphere, 0),
		Damping:   damping,
		TimeStep:  timeStep,
		SubSteps:  subSteps,
		Dt:        timeStep / float32(subSteps),
		Bounds:    bounds,
		PathExp:   pathExp,
		HisExp:    hisExp,
	}
}

func (e *Engine) AddSphere(s models.Sphere) {
	e.Objects = append(e.Objects, &s)
}

func (e *Engine) AddObstacle(s models.Sphere) {
	e.Obstacles = append(e.Obstacles, &s)
}

func (e *Engine) Update() {
	for i := 0; i < e.SubSteps; i++ {
		e.updateBodiesParallel()
		e.validateCollisionsParallel()
	}
}

func (e *Engine) CreateMesh() {
	e.Mesh = models.NewMesh(e.Bounds, e.MaxSize)
}

func (e *Engine) UpdateMesh() {
	for i := 0; i < len(e.Objects); i++ {
		e.Mesh.AddObject(e.Objects[i].Position, i)
	}

	for i := 0; i < len(e.Obstacles); i++ {
		e.Mesh.AddObstacle(e.Obstacles[i].Position, i)
	}
}

func (e *Engine) updateBodiesParallel() {
	var wg sync.WaitGroup

	for i := 0; i < len(e.Objects); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.Objects[i].ApplyForce(&e.Gravity)
			e.Objects[i].Update(e.Dt)
			CollisionBounds(e.Objects[i], e.Bounds, e.Damping)
		}(i)
	}

	wg.Wait()
}

func (e *Engine) validateCollisionsParallel() {
	e.UpdateMesh()

	verticalNumDivisions := 10
	horizontalNumDivisions := 10

	verticalDivisionSize := e.Mesh.Rows / verticalNumDivisions
	horizontalDivisionSize := e.Mesh.Columns / horizontalNumDivisions

	var wg sync.WaitGroup

	for sectVertical := 0; sectVertical <= verticalNumDivisions; sectVertical++ {
		for sectHorizontal := 0; sectHorizontal <= horizontalNumDivisions; sectHorizontal++ {
			startX := sectHorizontal * horizontalDivisionSize
			startY := sectVertical * verticalDivisionSize

			wg.Add(1)
			go func(startX, startY, horizontalDivisionSize, verticalDivisionSize int) {
				defer wg.Done()

				for i := startX; i < startX+horizontalDivisionSize; i++ {
					for j := startY; j < startY+verticalDivisionSize; j++ {
						objects, obstacles := e.Mesh.GetElementsAround(i, j)

						for k := 0; k < len(objects); k++ {
							for l := 0; l < len(obstacles); l++ {
								ValidateCollision(e.Objects[*objects[k]], e.Obstacles[*obstacles[l]])
							}

							for m := k + 1; m < len(objects); m++ {
								sA := e.Objects[*objects[k]]
								sB := e.Objects[*objects[m]]

								if sA != sB {
									ValidateCollision(sA, sB)
								}
							}
						}
					}
				}

			}(startX, startY, horizontalDivisionSize, verticalDivisionSize)
		}
	}

	wg.Wait()
	e.Mesh.Clear()
}

func (e *Engine) ExportCurrentState(borders []*gmath.Vector) {
	total := len(e.Objects) + len(e.Obstacles) + len(borders)

	content := export.GetExportHeader(total)
	e.PathExp.Write(content)

	itemCount := 0
	for itemCount < len(e.Objects) {
		content = export.GetExportPath(itemCount, e.Objects[itemCount])
		e.PathExp.Write(content)
		itemCount++
	}
	for j := 0; j < len(e.Obstacles); j++ {
		content = export.GetExportPath(itemCount, e.Obstacles[j])
		e.PathExp.Write(content)
		itemCount++
	}
	for j := 0; j < len(borders); j++ {
		content = export.GetExportPathBorders(itemCount, borders[j])
		e.PathExp.Write(content)
		itemCount++
	}
}

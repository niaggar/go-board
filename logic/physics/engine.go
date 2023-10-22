package physics

import (
	"go-board/export"
	"go-board/gmath"
	"go-board/models"
	"sync"
)

type Engine struct {
	Gravity      gmath.Vector
	Objects      []*models.Sphere
	Obstacles    []*models.Sphere
	Damping      float32
	TimeStep     float32
	SubSteps     int
	Dt           float32
	Bounds       gmath.Vector
	PathExp      *export.Exporter
	HisExp       *export.Exporter
	Mesh         *models.Mesh
	MaxSize      float32
	TotalObjects int
	Columns      int
	ColumnsSize  float32
	IsEnded      bool
}

func NewEngine(gravity, bounds gmath.Vector, damping, timeStep, columnsSize float32, subSteps, totalObjects, columns int, pathExp, hisExp *export.Exporter) Engine {
	return Engine{
		Gravity:      gravity,
		Objects:      make([]*models.Sphere, 0),
		Obstacles:    make([]*models.Sphere, 0),
		Damping:      damping,
		TimeStep:     timeStep,
		SubSteps:     subSteps,
		Dt:           timeStep / float32(subSteps),
		Bounds:       bounds,
		PathExp:      pathExp,
		HisExp:       hisExp,
		TotalObjects: totalObjects,
		Columns:      columns,
		ColumnsSize:  columnsSize,
	}
}

func (e *Engine) Update() {
	for i := 0; i < e.SubSteps; i++ {
		e.updateBodiesParallel()
		e.validateCollisionsParallel()
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
			CollisionBounds(e.Objects[i], e.Bounds, e.Damping, true)
		}(i)
	}

	wg.Wait()

	countInactiveElements := 0
	for i := 0; i < len(e.Objects); i++ {
		if !e.Objects[i].IsActive {
			countInactiveElements++
		}
	}

	if countInactiveElements == e.TotalObjects {
		e.IsEnded = true
	}
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
							if !e.Objects[*objects[k]].IsActive {
								continue
							}

							for l := 0; l < len(obstacles); l++ {
								if !e.Obstacles[*obstacles[l]].IsActive {
									continue
								}

								ValidateCollision(e.Objects[*objects[k]], e.Obstacles[*obstacles[l]])
							}

							for m := k + 1; m < len(objects); m++ {
								if !e.Objects[*objects[m]].IsActive {
									continue
								}

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

func (e *Engine) ExportHistogram() {
	finalCountByColumn := make([]int, e.Columns)

	for i := 0; i < len(e.Objects); i++ {
		pos := int(e.Objects[i].Position.X / e.ColumnsSize)
		finalCountByColumn[pos]++
	}

	content := export.GetExportHistogram(finalCountByColumn)
	e.HisExp.Write(content)
}

func (e *Engine) AddSphere(s models.Sphere) {
	e.Objects = append(e.Objects, &s)
}

func (e *Engine) AddObstacle(s models.Sphere) {
	e.Obstacles = append(e.Obstacles, &s)
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

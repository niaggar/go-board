package physics

import (
	"go-board/export"
	"go-board/gmath"
	"go-board/models"
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
	Mesh      models.Mesh
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
		e.updateBodies()
		e.validateCollisions()
	}
}

func (e *Engine) updateBodies() {
	for i := 0; i < len(e.Objects); i++ {
		e.Objects[i].ApplyForce(&e.Gravity)
		e.Objects[i].Update(e.Dt)
		CollisionBounds(e.Objects[i], e.Bounds)
	}
}

func (e *Engine) validateCollisions() {
	for i := 0; i < len(e.Objects); i++ {
		for j := i + 1; j < len(e.Objects); j++ {
			ValidateCollision(e.Objects[i], e.Objects[j])
		}
		for _, o := range e.Obstacles {
			ValidateCollision(e.Objects[i], o)
		}
	}
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

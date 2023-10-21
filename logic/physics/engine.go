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
	Damping   float64
	TimeStep  float64
	SubSteps  int
	Dt        float64
	Bounds    gmath.Vector
	Exporter  *export.Exporter
	Mesh      models.Mesh
}

func NewEngine(gravity, bounds gmath.Vector, damping, timeStep float64, subSteps int, exporter *export.Exporter) Engine {
	return Engine{
		Gravity:   gravity,
		Objects:   make([]*models.Sphere, 0),
		Obstacles: make([]*models.Sphere, 0),
		Damping:   damping,
		TimeStep:  timeStep,
		SubSteps:  subSteps,
		Dt:        timeStep / float64(subSteps),
		Bounds:    bounds,
		Exporter:  exporter,
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

func (e *Engine) ExportCurrentState() {
	total := len(e.Objects) + len(e.Obstacles)
	e.Exporter.StartFrame(total)

	for _, s := range e.Objects {
		e.Exporter.ExportElement(*s)
	}
	for _, s := range e.Obstacles {
		e.Exporter.ExportElement(*s)
	}
}

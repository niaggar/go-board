package logic

import (
	"go-board/export"
	"go-board/gmath"
	"go-board/models"
)

type Engine struct {
	Gravity  gmath.Vector
	Objects  []*models.Sphere
	Damping  float64
	TimeStep float64
	SubSteps int
	Dt       float64
	Bounds   gmath.Vector
	Exporter *export.Exporter
	Mesh     models.Mesh
}

func NewEngine(gravity, bounds gmath.Vector, damping, timeStep float64, subSteps int, exporter *export.Exporter) *Engine {
	return &Engine{
		Gravity:  gravity,
		Objects:  make([]*models.Sphere, 0),
		Damping:  damping,
		TimeStep: timeStep,
		SubSteps: subSteps,
		Dt:       timeStep / float64(subSteps),
		Bounds:   bounds,
		Exporter: exporter,
	}
}

func (e *Engine) AddSphere(s *models.Sphere) {
	e.Objects = append(e.Objects, s)
}

func (e *Engine) Update() {
	for i := 0; i < e.SubSteps; i++ {
		for _, s := range e.Objects {
			s.ApplyForce(&e.Gravity)
			s.Update(e.Dt)
			s.CollisionBounds(e.Bounds)
		}

		for _, sA := range e.Objects {
			for _, sB := range e.Objects {
				if sA != sB {
					sA.ValidateCollision(sB)
				}
			}
		}
	}
}

func (e *Engine) Export() {
	total := len(e.Objects)
	e.Exporter.StartFrame(total)

	for _, s := range e.Objects {
		e.Exporter.ExportElement(*s)
	}
}

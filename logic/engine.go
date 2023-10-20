package logic

import (
	"go-board/export"
	"go-board/gmath"
	"go-board/models"
)

type Engine struct {
	gravity  gmath.Vector
	objects  []*models.Sphere
	damping  float64
	timeStep float64
	subSteps int
	dt       float64
	bounds   gmath.Vector
	exporter *export.Exporter
}

func NewEngine(gravity, bounds gmath.Vector, damping, timeStep float64, subSteps int) *Engine {
	return &Engine{
		gravity:  gravity,
		objects:  make([]*models.Sphere, 0),
		damping:  damping,
		timeStep: timeStep,
		subSteps: subSteps,
		dt:       timeStep / float64(subSteps),
		bounds:   bounds,
		exporter: export.NewExporter("./export.dat"),
	}
}

func (e *Engine) AddSphere(s *models.Sphere) {
	e.objects = append(e.objects, s)
}

func (e *Engine) Update() {
	for i := 0; i < e.subSteps; i++ {
		for _, s := range e.objects {
			s.ApplyForce(&e.gravity)
			s.Update(e.dt)
			s.CollisionBounds(e.bounds)
		}

		for _, sA := range e.objects {
			for _, sB := range e.objects {
				if sA != sB {
					sA.ValidateCollision(sB)
				}
			}
		}
	}
}

func (e *Engine) Export() {
	total := len(e.objects)
	e.exporter.StartFrame(total)

	for _, s := range e.objects {
		e.exporter.ExportElement(*s)
	}
}

func (e *Engine) CreateFile() {
	e.exporter.CreateFile()
}

func (e *Engine) CloseFile() {
	e.exporter.CloseFile()
}

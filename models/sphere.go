package models

import (
	"go-board/gmath"
)

type Sphere struct {
	Id          int
	Damping     float32
	Radius      float32
	Mass        float32
	InverseMass float32
	Position    gmath.Vector
	Velocity    gmath.Vector
	Force       gmath.Vector
	Type        int
	CanCollide  bool
	IsActive    bool
}

func NewSphere(x, y, radius, mass, damping float32, t int) *Sphere {
	return &Sphere{
		Id:          gmath.GetRandomID(),
		Damping:     damping,
		Radius:      radius,
		Mass:        mass,
		InverseMass: 1 / mass,
		Position:    gmath.NewVector(x, y),
		Velocity:    gmath.NewVector(0, 0),
		Force:       gmath.NewVector(0, 0),
		Type:        t,
		IsActive:    true,
	}
}

func (s *Sphere) Update(dt float32) {
	if s.Type == STATIC || !s.IsActive {
		return
	}

	newVelocity := gmath.Scale(s.Force, s.InverseMass*dt)
	s.Velocity = gmath.Add(s.Velocity, newVelocity)

	newPosition := gmath.Scale(s.Velocity, dt)
	s.Position = gmath.Add(s.Position, newPosition)

	s.Force = gmath.NewVector(0, 0)
}

func (s *Sphere) ApplyForce(f *gmath.Vector) {
	if s.Type == STATIC || !s.IsActive {
		return
	}

	s.Force = gmath.Add(s.Force, *f)
}

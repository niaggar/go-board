package models

import (
	"go-board/gmath"
	"math"
)

type Sphere struct {
	Damping     float64
	Radius      float64
	Mass        float64
	InverseMass float64
	Position    gmath.Vector
	Velocity    gmath.Vector
	Force       gmath.Vector
	Type        int
}

func NewSphere(x, y, radius, mass, damping float64, t int) *Sphere {
	return &Sphere{
		Damping:     damping,
		Radius:      radius,
		Mass:        mass,
		InverseMass: 1 / mass,
		Position:    gmath.NewVector(x, y),
		Velocity:    gmath.NewVector(0, 0),
		Force:       gmath.NewVector(0, 0),
		Type:        t,
	}
}

func (s *Sphere) Update(dt float64) {
	newVelocity := gmath.Scale(&s.Force, s.InverseMass*dt)
	s.Velocity = *gmath.Add(&s.Velocity, newVelocity)

	newPosition := gmath.Scale(&s.Velocity, dt)
	s.Position = *gmath.Add(&s.Position, newPosition)

	s.Force = gmath.NewVector(0, 0)
}

func (s *Sphere) ApplyForce(f *gmath.Vector) {
	s.Force = *gmath.Add(&s.Force, f)
}

func (sA *Sphere) ValidateCollision(sB *Sphere) {
	if sA.Type == STATIC && sB.Type == STATIC {
		return
	}

	rmin := sA.Radius + sB.Radius
	rminSqu := rmin * rmin
	direction := gmath.Sub(&sA.Position, &sB.Position)
	distanceSqu := gmath.LengthSqu(direction)

	if distanceSqu < rminSqu {
		// Separate spheres
		norm := gmath.Normalice(direction)
		overlap := math.Abs(rmin - math.Sqrt(distanceSqu))
		separateSphere(sA, sB, *norm, overlap)

		// Resolve collision
		reducedMass := 1 / (sA.InverseMass + sB.InverseMass)
		relativeVelocity := gmath.Sub(&sA.Velocity, &sB.Velocity)
		normalVelocity := gmath.Dot(relativeVelocity, norm)

		if normalVelocity < 0 {
			return
		}

		e := math.Min(sA.Damping, sB.Damping)
		j := -(1 + e) * normalVelocity * reducedMass
		impulse := gmath.Scale(norm, j)

		if sA.Type == DYNAMIC {
			sA.Velocity = *gmath.Add(&sA.Velocity, gmath.Scale(impulse, sA.InverseMass))
		}
		if sB.Type == DYNAMIC {
			sB.Velocity = *gmath.Sub(&sB.Velocity, gmath.Scale(impulse, sB.InverseMass))
		}
	}
}

func (s *Sphere) CollisionBounds(bounds gmath.Vector) {
	possX := s.Position.X()
	possY := s.Position.Y()

	if possX-s.Radius < 0 {
		s.Position.SetX(s.Radius)
		s.Velocity.SetX(-s.Velocity.X() * s.Damping)
	} else if possX+s.Radius > bounds.X() {
		s.Position.SetX(bounds.X() - s.Radius)
		s.Velocity.SetX(-s.Velocity.X() * s.Damping)
	}

	if possY-s.Radius < 0 {
		s.Position.SetY(s.Radius)
		s.Velocity.SetY(-s.Velocity.Y() * s.Damping)
	} else if possY+s.Radius > bounds.Y() {
		s.Position.SetY(bounds.Y() - s.Radius)
		s.Velocity.SetY(-s.Velocity.Y() * s.Damping)
	}
}

func separateSphere(sA, sB *Sphere, norm gmath.Vector, overlap float64) {
	if sA.Type == STATIC {
		sB.Position = *gmath.Sub(&sA.Position, gmath.Scale(&norm, overlap))
	} else if sB.Type == STATIC {
		sA.Position = *gmath.Add(&sB.Position, gmath.Scale(&norm, overlap))
	} else {
		sA.Position = *gmath.Add(&sA.Position, gmath.Scale(&norm, overlap/2))
		sB.Position = *gmath.Sub(&sB.Position, gmath.Scale(&norm, overlap/2))
	}
}

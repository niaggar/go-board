package models

import (
	"go-board/bmath"
	"math"
)

type Sphere struct {
	Damping     float64
	Radius      float64
	Mass        float64
	InverseMass float64
	Position    bmath.Vector
	Velocity    bmath.Vector
	Force       bmath.Vector
	Type        int
}

func NewSphere(x, y, radius, mass, damping float64, t int) *Sphere {
	return &Sphere{
		Damping:     damping,
		Radius:      radius,
		Mass:        mass,
		InverseMass: 1 / mass,
		Position:    bmath.NewVector(x, y),
		Velocity:    bmath.NewVector(0, 0),
		Force:       bmath.NewVector(0, 0),
		Type:        t,
	}
}

func (s *Sphere) Update(dt float64) {
	newVelocity := bmath.Scale(&s.Force, s.InverseMass*dt)
	s.Velocity = *bmath.Add(&s.Velocity, newVelocity)

	newPosition := bmath.Scale(&s.Velocity, dt)
	s.Position = *bmath.Add(&s.Position, newPosition)

	s.Force = bmath.NewVector(0, 0)
}

func (s *Sphere) ApplyForce(f *bmath.Vector) {
	s.Force = *bmath.Add(&s.Force, f)
}

func (sA *Sphere) ValidateCollision(sB *Sphere) {
	if sA.Type == STATIC && sB.Type == STATIC {
		return
	}

	rmin := sA.Radius + sB.Radius
	rminSqu := rmin * rmin
	direction := bmath.Sub(&sA.Position, &sB.Position)
	distanceSqu := bmath.LengthSqu(direction)

	if distanceSqu < rminSqu {
		// Separate spheres
		norm := bmath.Normalice(direction)
		overlap := math.Abs(rmin - math.Sqrt(distanceSqu))
		separateSphere(sA, sB, *norm, overlap)

		// Resolve collision
		reducedMass := 1 / (sA.InverseMass + sB.InverseMass)
		relativeVelocity := bmath.Sub(&sA.Velocity, &sB.Velocity)
		normalVelocity := bmath.Dot(relativeVelocity, norm)

		if normalVelocity < 0 {
			return
		}

		e := math.Min(sA.Damping, sB.Damping)
		j := -(1 + e) * normalVelocity * reducedMass
		impulse := bmath.Scale(norm, j)

		if sA.Type == DYNAMIC {
			sA.Velocity = *bmath.Add(&sA.Velocity, bmath.Scale(impulse, sA.InverseMass))
		}
		if sB.Type == DYNAMIC {
			sB.Velocity = *bmath.Sub(&sB.Velocity, bmath.Scale(impulse, sB.InverseMass))
		}
	}
}

func (s *Sphere) CollisionBounds(bounds bmath.Vector) {
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

func separateSphere(sA, sB *Sphere, norm bmath.Vector, overlap float64) {
	if sA.Type == STATIC {
		sB.Position = *bmath.Sub(&sA.Position, bmath.Scale(&norm, overlap))
	} else if sB.Type == STATIC {
		sA.Position = *bmath.Add(&sB.Position, bmath.Scale(&norm, overlap))
	} else {
		sA.Position = *bmath.Add(&sA.Position, bmath.Scale(&norm, overlap/2))
		sB.Position = *bmath.Sub(&sB.Position, bmath.Scale(&norm, overlap/2))
	}
}

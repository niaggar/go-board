package physics

import (
	"go-board/gmath"
	"go-board/models"
	"math"
)

func CollisionBounds(s *models.Sphere, bounds gmath.Vector, damping float32) {
	if s.Position.X-s.Radius < 0 {
		s.Position.X = s.Radius
		s.Velocity.X = -s.Velocity.X * damping
	} else if s.Position.X+s.Radius > bounds.X {
		s.Position.X = bounds.X - s.Radius
		s.Velocity.X = -s.Velocity.X * damping
	}

	if s.Position.Y-s.Radius < 0 {
		s.Position.Y = s.Radius
		s.Velocity.Y = -s.Velocity.Y * damping
	} else if s.Position.Y+s.Radius > bounds.Y {
		s.Position.Y = bounds.Y - s.Radius
		s.Velocity.Y = -s.Velocity.Y * damping
	}
}

func ValidateCollision(sA, sB *models.Sphere) {
	if sA.Type == models.STATIC && sB.Type == models.STATIC {
		return
	}

	rmin := sA.Radius + sB.Radius
	rminSqu := rmin * rmin
	direction := gmath.Sub(sA.Position, sB.Position)
	distanceSqu := gmath.LengthSqu(direction)

	if distanceSqu < rminSqu {
		norm := gmath.Normalice(direction)
		overlap := float32(math.Sqrt(float64(distanceSqu)) - float64(rmin))

		separateSphere(sA, sB, norm, overlap)
		resolveCollision(sA, sB, norm)
	}
}

func separateSphere(sA, sB *models.Sphere, norm gmath.Vector, overlap float32) {
	if sA.Type == models.STATIC {
		sB.Position = gmath.Add(sB.Position, gmath.Scale(norm, overlap))
	} else if sB.Type == models.STATIC {
		sA.Position = gmath.Sub(sA.Position, gmath.Scale(norm, overlap))
	} else {
		sA.Position = gmath.Sub(sA.Position, gmath.Scale(norm, overlap/2))
		sB.Position = gmath.Add(sB.Position, gmath.Scale(norm, overlap/2))
	}
}

func resolveCollision(sA, sB *models.Sphere, norm gmath.Vector) {
	if sA.Type == models.STATIC && sB.Type == models.STATIC {
		return
	}

	reducedMass := 1.0 / (sA.InverseMass + sB.InverseMass)
	vA := sA.Velocity
	vB := sB.Velocity

	relativeVelocity := gmath.Sub(vA, vB)
	velocityAlongNormal := gmath.Dot(norm, relativeVelocity)
	if velocityAlongNormal > 0 {
		return
	}

	e := math.Min(float64(sA.Damping), float64(sB.Damping))
	j := (-velocityAlongNormal * (float32(1) + float32(e))) * reducedMass

	impulse := gmath.Scale(norm, j)

	if sA.Type == models.DYNAMIC {
		sA.Velocity = gmath.Add(sA.Velocity, gmath.Scale(impulse, sA.InverseMass))
	}
	if sB.Type == models.DYNAMIC {
		sB.Velocity = gmath.Sub(sB.Velocity, gmath.Scale(impulse, sB.InverseMass))
	}
}

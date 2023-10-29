package physics

import (
	"go-board/models"
	"go-board/utils/gmath"
	"math"
)

func CollisionBounds(s *models.Ball, boardProps *models.BoardProps, bounds *models.Bounds, detectFloor bool) {
	if s.Position.X-s.Radius < 0 {
		s.Position.X = s.Radius
		s.Velocity.X = -s.Velocity.X * boardProps.Damping
	} else if s.Position.X+s.Radius > bounds.Width {
		s.Position.X = bounds.Width - s.Radius
		s.Velocity.X = -s.Velocity.X * boardProps.Damping
	}

	if s.Position.Y-s.Radius < 0 {
		s.Position.Y = s.Radius
		s.Velocity.Y = -s.Velocity.Y * boardProps.Damping

		if detectFloor {
			s.Active = false
		}
	} else if s.Position.Y+s.Radius > bounds.Height {
		s.Position.Y = bounds.Height - s.Radius
		s.Velocity.Y = -s.Velocity.Y * boardProps.Damping
	}
}

func ValidateCollision(sA, sB *models.Ball) {
	if sA.Static && sB.Static {
		return
	}

	rMin := sA.Radius + sB.Radius
	rMinSqr := rMin * rMin
	direction := gmath.Sub(sA.Position, sB.Position)
	distanceSqu := gmath.LengthSqu(direction)

	if distanceSqu < rMinSqr {
		norm := gmath.Normalice(direction)
		overlap := float32(math.Sqrt(float64(distanceSqu)) - float64(rMin))

		separateSphere(sA, sB, norm, overlap)
		resolveCollision(sA, sB, norm)
	}
}

func separateSphere(sA, sB *models.Ball, norm gmath.Vector, overlap float32) {
	if sA.Static {
		sB.Position = gmath.Add(sB.Position, gmath.Scale(norm, overlap))
	} else if sB.Static {
		sA.Position = gmath.Sub(sA.Position, gmath.Scale(norm, overlap))
	} else {
		sA.Position = gmath.Sub(sA.Position, gmath.Scale(norm, overlap/2))
		sB.Position = gmath.Add(sB.Position, gmath.Scale(norm, overlap/2))
	}
}

func resolveCollision(sA, sB *models.Ball, norm gmath.Vector) {
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

	if !sA.Static {
		sA.Velocity = gmath.Add(sA.Velocity, gmath.Scale(impulse, sA.InverseMass))
	}
	if !sB.Static {
		sB.Velocity = gmath.Sub(sB.Velocity, gmath.Scale(impulse, sB.InverseMass))
	}
}

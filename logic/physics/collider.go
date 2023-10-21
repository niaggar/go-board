package physics

import (
	"go-board/gmath"
	"go-board/models"
	"math"
)

func CollisionBounds(s *models.Sphere, bounds gmath.Vector) {
	if s.Position.X-s.Radius < 0 {
		s.Position.X = s.Radius
		s.Velocity.X = -s.Velocity.X * s.Damping
	} else if s.Position.X+s.Radius > bounds.X {
		s.Position.X = bounds.X - s.Radius
		s.Velocity.X = -s.Velocity.X * s.Damping
	}

	if s.Position.Y-s.Radius < 0 {
		s.Position.Y = s.Radius
		s.Velocity.Y = -s.Velocity.Y * s.Damping
	} else if s.Position.Y+s.Radius > bounds.Y {
		s.Position.Y = bounds.Y - s.Radius
		s.Velocity.Y = -s.Velocity.Y * s.Damping
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
		resolveCollision(sA, sB, norm, overlap)
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

func resolveCollision(sA, sB *models.Sphere, norm gmath.Vector, overlap float32) {
	tangent := gmath.Vector{X: -norm.Y, Y: norm.X}
	totalMass := sA.Mass + sB.Mass

	// Project velocities onto the collision normal and tangent
	vAn := gmath.Dot(norm, sA.Velocity)
	vAt := gmath.Dot(tangent, sA.Velocity)
	vBn := gmath.Dot(norm, sB.Velocity)
	vBt := gmath.Dot(tangent, sB.Velocity)

	// Calculate the new normal velocities
	vAnNew := (vAn*(sA.Mass-sB.Mass) + 2*sB.Mass*vBn) / totalMass
	vBnNew := (vBn*(sB.Mass-sA.Mass) + 2*sA.Mass*vAn) / totalMass

	// Calculate the new tangent velocities
	vAtNew := vAt
	vBtNew := vBt

	// Calculate the new normal and tangent vectors
	damping := float32(math.Min(float64(sA.Damping), float64(sB.Damping)))

	vAnNewNorm := gmath.Scale(norm, vAnNew*damping)
	vAnNewTang := gmath.Scale(tangent, vAtNew)
	vBnNewNorm := gmath.Scale(norm, vBnNew*damping)
	vBnNewTang := gmath.Scale(tangent, vBtNew)

	// Calculate the new velocities
	sA.Velocity = gmath.Add(vAnNewNorm, vAnNewTang)
	sB.Velocity = gmath.Add(vBnNewNorm, vBnNewTang)
}

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
	direction := gmath.Sub(&sA.Position, &sB.Position)
	distanceSqu := gmath.LengthSqu(direction)

	if distanceSqu < rminSqu {
		norm := gmath.Normalice(direction)
		overlap := math.Abs(rmin - math.Sqrt(distanceSqu))

		separateSphere(sA, sB, *norm, overlap)
		resolveCollision(sA, sB, *norm, overlap)
	}
}

func separateSphere(sA, sB *models.Sphere, norm gmath.Vector, overlap float64) {
	if sA.Type == models.STATIC {
		sB.Position = *gmath.Sub(&sA.Position, gmath.Scale(&norm, overlap))
	} else if sB.Type == models.STATIC {
		sA.Position = *gmath.Add(&sB.Position, gmath.Scale(&norm, overlap))
	} else {
		sA.Position = *gmath.Add(&sA.Position, gmath.Scale(&norm, overlap/2))
		sB.Position = *gmath.Sub(&sB.Position, gmath.Scale(&norm, overlap/2))
	}
}

func resolveCollision(sA, sB *models.Sphere, norm gmath.Vector, overlap float64) {
	dx := sA.Position.X - sB.Position.X
	dy := sA.Position.Y - sB.Position.Y
	impact_angle := math.Atan(dy / dx)

	vxa := sA.Velocity.X
	vya := sA.Velocity.Y
	aplha_0 := sB.Damping

	v_tangen := -vxa*math.Sin(impact_angle) + vya*math.Cos(impact_angle)
	v_radial := -aplha_0 * (vxa*math.Cos(impact_angle) + vya*math.Sin(impact_angle))

	vx_new := v_radial*math.Cos(impact_angle) - v_tangen*math.Sin(impact_angle)
	vy_new := v_radial*math.Sin(impact_angle) + v_tangen*math.Cos(impact_angle)

	new_x := (sum_radius)*math.cos(impact_angle) + position_o.x
	new_y := (sum_radius)*math.sin(impact_angle) + position_o.y

	esphere1.set_position(Vect2(new_x, new_y))
	esphere1.set_velocity(Vect2(vx_new, vy_new))
}

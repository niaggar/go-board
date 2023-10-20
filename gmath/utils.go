package gmath

import "math"

func Normalice(v *Vector) *Vector {
	if v.x == 0 && v.y == 0 {
		return v
	}

	n := v
	length := Length(n)

	return Scale(n, 1/length)
}

func Dot(v1, v2 *Vector) float64 {
	return v1.x*v2.x + v1.y*v2.y
}

func Cross(v1, v2 *Vector) float64 {
	return v1.x*v2.y - v1.y*v2.x
}

func Add(v1, v2 *Vector) *Vector {
	return &Vector{v1.x + v2.x, v1.y + v2.y}
}

func Sub(v1, v2 *Vector) *Vector {
	return &Vector{v1.x - v2.x, v1.y - v2.y}
}

func Scale(v *Vector, s float64) *Vector {
	return &Vector{v.x * s, v.y * s}
}

func Length(v *Vector) float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func LengthSqu(v *Vector) float64 {
	return v.x*v.x + v.y*v.y
}

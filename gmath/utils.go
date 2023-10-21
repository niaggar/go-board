package gmath

import "math"

func Normalice(v *Vector) *Vector {
	if v.X == 0 && v.Y == 0 {
		return v
	}

	n := v
	length := Length(n)

	return Scale(n, 1/length)
}

func Dot(v1, v2 *Vector) float64 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func Cross(v1, v2 *Vector) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

func Add(v1, v2 *Vector) *Vector {
	return &Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func Sub(v1, v2 *Vector) *Vector {
	return &Vector{v1.X - v2.X, v1.Y - v2.Y}
}

func Scale(v *Vector, s float64) *Vector {
	return &Vector{v.X * s, v.Y * s}
}

func Length(v *Vector) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func LengthSqu(v *Vector) float64 {
	return v.X*v.X + v.Y*v.Y
}

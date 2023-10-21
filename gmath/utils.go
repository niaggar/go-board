package gmath

import "math"

func Normalice(v Vector) Vector {
	if v.X == 0 && v.Y == 0 {
		return v
	}

	length := Length(v)
	return Scale(v, 1/length)
}

func Dot(v1, v2 Vector) float32 {
	return v1.X*v2.X + v1.Y*v2.Y
}

func Cross(v1, v2 Vector) float32 {
	return v1.X*v2.Y - v1.Y*v2.X
}

func Add(v1, v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func Sub(v1, v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y}
}

func Scale(v Vector, s float32) Vector {
	return Vector{v.X * s, v.Y * s}
}

func Length(v Vector) float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func LengthSqu(v Vector) float32 {
	return v.X*v.X + v.Y*v.Y
}

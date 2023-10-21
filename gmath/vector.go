package gmath

type Vector struct {
	X float32
	Y float32
}

func NewVector(x, y float32) Vector {
	return Vector{x, y}
}

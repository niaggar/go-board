package gmath

type Vector struct {
	x float64
	y float64
}

func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

func (v *Vector) X() float64 {
	return v.x
}

func (v *Vector) Y() float64 {
	return v.y
}

func (v *Vector) SetX(x float64) {
	v.x = x
}

func (v *Vector) SetY(y float64) {
	v.y = y
}

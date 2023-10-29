package models

import (
	"go-board/utils/gmath"
)

type BallProps struct {
	Damping     float32
	Radius      float32
	Mass        float32
	InverseMass float32
}

type Ball struct {
	BallProps
	Position gmath.Vector
	Velocity gmath.Vector
	Force    gmath.Vector
	Active   bool // Seguir actualizando sus posiciones y detectando colisiones
	Static   bool // Objeto no se mueve
	Obstacle bool
}

func NewBall(x, y float32, active, static, obstacle bool, props BallProps) *Ball {
	return &Ball{
		BallProps: props,
		Position:  gmath.NewVector(x, y),
		Velocity:  gmath.NewVector(0, 0),
		Force:     gmath.NewVector(0, 0),
		Active:    active,
		Static:    static,
		Obstacle:  obstacle,
	}
}

func (ball *Ball) Update(dt float32) {
	if !ball.Active {
		return
	}

	newVelocity := gmath.Scale(ball.Force, ball.InverseMass*dt)
	ball.Velocity = gmath.Add(ball.Velocity, newVelocity)

	newPosition := gmath.Scale(ball.Velocity, dt)
	ball.Position = gmath.Add(ball.Position, newPosition)

	ball.Force = gmath.NewVector(0, 0)
}

func (ball *Ball) ApplyForce(f *gmath.Vector) {
	if !ball.Active {
		return
	}

	ball.Force = gmath.Add(ball.Force, *f)
}

func (ball *Ball) SetVelocity(vx, vy float32) {
	ball.Velocity = gmath.Vector{X: vx, Y: vy}
}

func (ball *Ball) SetPosition(x, y float32) {
	ball.Position = gmath.Vector{X: x, Y: y}
}

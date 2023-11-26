package models

import (
	"go-board/utils/gmath"
	"math"
)

type ObstacleProps struct {
	XAmplitude, YAmplitude float32
	XFrequency, YFrequency float32
	CurrentTime            float32
}

type BallProps struct {
	Damping     float32
	Radius      float32
	Mass        float32
	InverseMass float32
}

type Ball struct {
	BallProps
	ObstacleProps
	Position gmath.Vector
	Velocity gmath.Vector
	Force    gmath.Vector
	Active   bool
	Static   bool
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

	if ball.Obstacle {
		ball.updateObstacle(dt)
	} else {
		ball.updateBall(dt)
	}
}

func (ball *Ball) updateBall(dt float32) {
	newVelocity := gmath.Scale(ball.Force, ball.InverseMass*dt)
	ball.Velocity = gmath.Add(ball.Velocity, newVelocity)

	newPosition := gmath.Scale(ball.Velocity, dt)
	ball.Position = gmath.Add(ball.Position, newPosition)

	ball.Force = gmath.NewVector(0, 0)
}

func (ball *Ball) updateObstacle(dt float32) {
	ball.CurrentTime += dt

	xFunction := SinusoidalFunction(ball.CurrentTime, ball.XFrequency, ball.XAmplitude)
	yFunction := SinusoidalFunction(ball.CurrentTime, ball.YFrequency, ball.YAmplitude)

	xNew := ball.Position.X + xFunction
	yNew := ball.Position.Y + yFunction

	ball.Position = gmath.Vector{X: xNew, Y: yNew}
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

func SinusoidalFunction(x, w, a float32) float32 {
	return a * float32(math.Sin(float64(x*w)))
}

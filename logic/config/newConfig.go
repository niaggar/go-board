package config

import "go-board/gmath"

type object struct {
	Radius    float64
	Damping   float64
	Mass      float64
	Collision bool
	Position  gmath.Vector
	Velocity  gmath.Vector
}

type minMax struct {
	Min float64
	Max float64
}

type board struct {
	ExportRoute  string
	GridSize     float64
	RowNumber    int
	ColumnNumber int
	Damping      float64
	Gravity      gmath.Vector
	TimeStep     float64
	SubSteps     int
	MaxTime      struct {
		MaxTime     float64
		TouchGround bool
	}
}

type NewConfig struct {
	Board       board
	CreateBalls struct {
		Collisions bool
		Positions  []object
		Creation   struct {
			Enabled  bool
			Count    int
			Damping  float64
			Mass     float64
			Radius   minMax
			Position struct {
				X minMax
				Y minMax
			}
			Velocity struct {
				X minMax
				Y minMax
			}
		}
	}
	CreateObstacles struct {
		Positions []object
		Creation  struct {
			Enabled   bool
			Damping   float64
			Mass      float64
			Radius    minMax
			Direction int
			YOffset   float64
		}
	}
}

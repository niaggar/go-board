package config

import "go-board/gmath"

type object struct {
	Radius    float32
	Damping   float32
	Mass      float32
	Collision bool
	Position  gmath.Vector
	Velocity  gmath.Vector
}

type minMax struct {
	Min float32
	Max float32
}

type export struct {
	Enabled bool
	Path    string
}

type board struct {
	GridSize     float32
	RowNumber    int
	ColumnNumber int
	Damping      float32
	Gravity      gmath.Vector
	TimeStep     float32
	SubSteps     int
	MaxTime      struct {
		MaxTime     float32
		TouchGround bool
	}
}

type NewConfig struct {
	Board       board
	ExportPath  export
	ExportHisto export
	CreateBalls struct {
		Collisions bool
		Positions  []object
		Creation   struct {
			Enabled  bool
			Count    int
			Damping  float32
			Mass     float32
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
			Damping   float32
			Mass      float32
			Radius    minMax
			Direction int
			YOffset   float32
		}
	}
}

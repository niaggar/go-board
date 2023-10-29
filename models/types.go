package models

import "go-board/utils/gmath"

type MinMax struct {
	Min float32
	Max float32
}

type Bounds struct {
	Width  float32
	Height float32
}

type PointMinMax struct {
	X MinMax
	Y MinMax
}

type SimulationProps struct {
	TimeStep   float32
	SubSteps   int
	Gravity    gmath.Vector
	MaxSize    float32
	Collisions bool
	Stop       struct {
		MaxTime     float32
		TouchGround bool
	}
}

type BoardProps struct {
	CellSize     float32
	RowNumber    int
	ColumnNumber int
	Damping      float32
}

type ExportProps struct {
	Active bool
	Route  string
}

type ExperimentProps struct {
	Name            string
	FrameRate       int
	Executions      int
	CreationFrame   int
	CreationNum     int
	ExportPaths     ExportProps
	ExportHistogram ExportProps
}

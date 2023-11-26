package models

type CreationBalls struct {
	Existing []Ball
	Creation struct {
		Enabled  bool
		Count    int
		Damping  float32
		Mass     float32
		Radius   MinMax
		Position PointMinMax
		Velocity PointMinMax
	}
}

type CreationObstacles struct {
	Existing []Ball
	Creation struct {
		Enabled bool
		Damping float32
		Mass    float32
		Radius  struct {
			Range     MinMax
			Direction int
		}
		Movement struct {
			Enabled  bool
			Range    MinMax
			TimeStep MinMax
		}
		XAmplitude, YAmplitude float32
		XFrequency, YFrequency float32
	}
}

type ConfigurationFile struct {
	Experiment ExperimentProps
	Board      BoardProps
	Simulation SimulationProps
	Balls      CreationBalls
	Obstacles  CreationObstacles
}

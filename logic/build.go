package logic

import (
	"go-board/models"
	"go-board/utils/gmath"
)

func BuildSpheres(props *models.CreationBalls, bounds *models.Bounds) []*models.Ball {
	autoCreation := props.Creation.Enabled
	creationProps := props.Creation

	balls := make([]*models.Ball, 0)

	if autoCreation {
		ballProps := models.BallProps{
			Damping:     creationProps.Damping,
			Mass:        creationProps.Mass,
			InverseMass: 1 / creationProps.Mass,
		}

		for i := 0; i < props.Creation.Count; i++ {
			ballProps.Radius = gmath.GetRandomFloat(creationProps.Radius.Min, creationProps.Radius.Max)

			xPos := bounds.Width/2 + gmath.GetRandomFloat(creationProps.Position.X.Min, creationProps.Position.X.Max)
			yPos := bounds.Height + gmath.GetRandomFloat(creationProps.Position.Y.Min, creationProps.Position.Y.Max)

			ball := models.NewBall(xPos, yPos, true, false, false, ballProps)
			balls = append(balls, ball)
		}
	}

	return balls
}

func BuildObstacles(props *models.CreationObstacles, board *models.BoardProps, bounds *models.Bounds) []*models.Ball {
	autoCreation := props.Creation.Enabled
	creationProps := props.Creation

	obstacles := make([]*models.Ball, 0)

	if autoCreation {
		radiusRange := creationProps.Radius.Range
		ballProps := models.BallProps{
			Damping:     creationProps.Damping,
			Mass:        creationProps.Mass,
			InverseMass: 1 / creationProps.Mass,
		}
		obsProps := models.ObstacleProps{
			XAmplitude: creationProps.XAmplitude,
			YAmplitude: creationProps.YAmplitude,
			XFrequency: creationProps.XFrequency,
			YFrequency: creationProps.YFrequency,
		}

		xOffset := bounds.Width / float32(board.ColumnNumber+1)
		yOffset := bounds.Height / float32(board.RowNumber+1)

		radiusDiff := radiusRange.Max - radiusRange.Min

		for i := 0; i < board.RowNumber; i++ {
			if creationProps.Radius.Direction == 0 {
				r := radiusRange.Min + (radiusDiff/float32(board.RowNumber))*float32(i)
				ballProps.Radius = r
			}

			for j := 0; j < board.ColumnNumber+2; j++ {
				if creationProps.Radius.Direction == 1 {
					r := radiusRange.Min + (radiusDiff/float32(board.ColumnNumber))*float32(j)
					ballProps.Radius = r
				}

				var xPos float32
				var yPos float32

				if i%2 == 0 {
					xPos = xOffset * float32(j)
					yPos = yOffset * float32(i+1)
				} else {
					if j >= board.ColumnNumber-1 {
						continue
					} else {
						xPos = xOffset*float32(j+1) + xOffset/2
						yPos = yOffset * float32(i+1)
					}
				}

				obstacle := models.NewBall(xPos, yPos, true, true, true, ballProps)
				obstacle.ObstacleProps = obsProps
				obstacles = append(obstacles, obstacle)
			}
		}
	}

	return obstacles
}

func BuildBorders(bounds *models.Bounds) []*gmath.Vector {
	borders := make([]*gmath.Vector, 0)

	for borderDivision := float32(0); borderDivision < bounds.Height; borderDivision += (0.5 * 2) {
		v1 := gmath.NewVector(0, borderDivision)
		v2 := gmath.NewVector(bounds.Width, borderDivision)

		borders = append(borders, &v1, &v2)
	}
	for borderDivision := float32(0); borderDivision < bounds.Width; borderDivision += (0.5 * 2) {
		v1 := gmath.NewVector(borderDivision, 0)
		v2 := gmath.NewVector(borderDivision, bounds.Height)

		borders = append(borders, &v1, &v2)
	}

	return borders
}

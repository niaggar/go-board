package bmath

import "math/rand"

func GetRandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func GetRandomFloat(min float64, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func GetRandomBool() bool {
	return rand.Intn(2) == 1
}

func GetRandomOnPoint(point float64, delta float64) float64 {
	return GetRandomFloat(point-delta, point+delta)
}

package gmath

import (
	"math/rand"
)

func GetRandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

func GetRandomFloat(min float32, max float32) float32 {
	return rand.Float32()*(max-min) + min
}

func GetRandomBool() bool {
	return rand.Intn(2) == 1
}

func GetRandomOnPoint(point float32, delta float32) float32 {
	return GetRandomFloat(point-delta, point+delta)
}

package gmath

import (
	"math/rand"
	"time"
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

func GetRandomID() int {
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(10000)
}

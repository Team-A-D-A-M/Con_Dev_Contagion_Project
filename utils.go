package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func floorDiv(a, b float32) int {
	return int(math.Floor(float64(a) / float64(b)))
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func limitVector(v rl.Vector2, max float32) rl.Vector2 {
	length := v.Length()

	if length > max {
		return rl.Vector2Scale(rl.Vector2Normalize(v), max)

	}

	return v
}

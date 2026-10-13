package main

import (
	"math"
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var sliderSeparation float32 = 25

func initGui() {
	var rect rl.Rectangle = rl.Rectangle{
		X:      40,
		Y:      0,
		Width:  200,
		Height: 20,
	}

	// max speed
	drawSlider(&maxSpeed, 1000, "Max Speed", &rect)

	// max force
	drawSlider(&maxForce, 200, "Max Force", &rect)

	// detection range
	drawSlider(&detectionRange, 200, "Detection Range", &rect)

	// desired separation
	drawSlider(&desiredSeparation, 200, "Desired Separation", &rect)

	//
	// cohesion value
	drawSlider(&cohesionValue, 1000, "Cohesion", &rect)

	// separation value
	drawSlider(&separationValue, 1000, "Separation", &rect)

	// alignment value
	drawSlider(&alignmentValue, 300, "Alignment", &rect)
}

func drawSlider(val *float32, max float32, label string, rect *rl.Rectangle) {
	rect.Y += float32(sliderSeparation)
	s := strconv.FormatFloat(roundFloat(float64(*val), 1), 'f', -1, 64)
	rg.Slider(*rect, s, label, val, 0, max)
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

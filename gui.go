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

	// boids amount
	s := strconv.Itoa(int(boidsAmount))
	tempVal := float32(boidsAmount)
	rg.Slider(rect, s, "Amount of Boids", &tempVal, 0, 1000)
	if int32(tempVal) != boidsAmount {
		setBoidsAmount(int32(tempVal))
	}

	// max speed
	rect.Y += sliderSeparation
	drawSlider(&maxSpeed, 1000, "Max Speed", &rect)

	// max force
	rect.Y += sliderSeparation
	drawSlider(&maxForce, 200, "Max Force", &rect)

	// detection range
	rect.Y += sliderSeparation
	drawSlider(&detectionRange, 200, "Detection Range", &rect)

	// desired separation
	rect.Y += sliderSeparation
	drawSlider(&desiredSeparation, 200, "Desired Separation", &rect)

	//
	// cohesion value
	rect.Y += sliderSeparation
	drawSlider(&cohesionValue, 1000, "Cohesion", &rect)

	// separation value
	rect.Y += sliderSeparation
	drawSlider(&separationValue, 1000, "Separation", &rect)

	// alignment value
	rect.Y += sliderSeparation
	drawSlider(&alignmentValue, 300, "Alignment", &rect)

	// uses grid
	rect.Y += sliderSeparation
	rg.CheckBox(rect, "Use Grid Optimization", &doGridOptimize)
}

func drawSlider(val *float32, max float32, label string, rect *rl.Rectangle) {
	s := strconv.FormatFloat(roundFloat(float64(*val), 1), 'f', -1, 64)
	rg.Slider(*rect, s, label, val, 0, max)
}

func drawSliderInt32(val *int32, max float32, label string, rect *rl.Rectangle) {
	s := string(*val)
	tempVal := float32(*val)
	rg.Slider(*rect, s, label, &tempVal, 0, max)
	*val = int32(tempVal)
}

func drawGrid() {
	// x
	total := int32(detectionRange)
	for total < int32(bounds.Width) {
		rl.DrawLine(total, 0, total, bounds.ToInt32().Height, rl.LightGray)
		total += int32(detectionRange)
	}

	// y
	total = int32(detectionRange)
	for total < int32(bounds.Height) {
		rl.DrawLine(0, total, bounds.ToInt32().Width, total, rl.LightGray)
		total += int32(detectionRange)
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

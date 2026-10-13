package ui

import (
	"boids/simulation"
	"strconv"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func InitGui() {
	var sliderSeparation float32 = 25.0
	var rect rl.Rectangle = rl.Rectangle{
		X:      40,
		Y:      0,
		Width:  200,
		Height: 20,
	}

	// boids amount
	s := strconv.Itoa(int(simulation.BoidsAmount))
	tempVal := float32(simulation.BoidsAmount)
	rg.Slider(rect, s, "Amount of Boids", &tempVal, 0, 1000)
	if int32(tempVal) != simulation.BoidsAmount {
		//main.SetBoidsAmount(int32(tempVal))
	}

	// max speed
	rect.Y += sliderSeparation
	drawSlider(&simulation.MaxSpeed, 1000, "Max Speed", &rect)

	// max force
	rect.Y += sliderSeparation
	drawSlider(&simulation.MaxForce, 200, "Max Force", &rect)

	// detection range
	rect.Y += sliderSeparation
	drawSlider(&simulation.DetectionRange, 200, "Detection Range", &rect)

	// desired separation
	rect.Y += sliderSeparation
	drawSlider(&simulation.DesiredSeparation, 200, "Desired Separation", &rect)

	//
	// cohesion value
	rect.Y += sliderSeparation
	drawSlider(&simulation.CohesionValue, 1000, "Cohesion", &rect)

	// separation value
	rect.Y += sliderSeparation
	drawSlider(&simulation.SeparationValue, 1000, "Separation", &rect)

	// alignment value
	rect.Y += sliderSeparation
	drawSlider(&simulation.AlignmentValue, 300, "Alignment", &rect)

	// uses grid
	rect.Y += sliderSeparation
	rg.CheckBox(rect, "Use Grid Optimization", &simulation.DoGridOptimize)
}

// draw parameter slider
func drawSlider(val *float32, max float32, label string, rect *rl.Rectangle) {
	s := strconv.FormatFloat(simulation.RoundFloat(float64(*val), 1), 'f', -1, 64)
	rg.Slider(*rect, s, label, val, 0, max)
}

// draw parameter slider for int
func drawSliderInt32(val *int32, max float32, label string, rect *rl.Rectangle) {
	s := string(*val)
	tempVal := float32(*val)
	rg.Slider(*rect, s, label, &tempVal, 0, max)
	*val = int32(tempVal)
}

// draw optimization grid
func DrawGrid() {
	// x
	total := int32(simulation.DetectionRange)
	for total < int32(simulation.Bounds.Width) {
		rl.DrawLine(total, 0, total, simulation.Bounds.ToInt32().Height, rl.LightGray)
		total += int32(simulation.DetectionRange)
	}
	// y
	total = int32(simulation.DetectionRange)
	for total < int32(simulation.Bounds.Height) {
		rl.DrawLine(0, total, simulation.Bounds.ToInt32().Width, total, rl.LightGray)
		total += int32(simulation.DetectionRange)
	}
}

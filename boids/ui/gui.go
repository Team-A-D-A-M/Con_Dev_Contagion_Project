// ui package initialises and draws all UI elements, including menu controls, boids and grid cells
package ui

import (
	"boids/simulation"
	"fmt"

	rg "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const SidebarWidth int32 = 280
const left float32 = 18
const sliderWidth float32 = 244

// InitGui draws the fixed sidebar and the seven parameter sliders.
func InitGui() {
	// Sidebar background.
	rl.DrawRectangle(
		0,
		0,
		SidebarWidth,
		simulation.ScreenY,
		rl.Color{R: 32, G: 43, B: 58, A: 255},
	)

	// Sidebar title.
	rl.DrawText("BOIDS", 18, 18, 26, rl.RayWhite)
	rl.DrawText("SIMULATION CONTROLS", 18, 51, 12, rl.LightGray)

	// Divider.
	rl.DrawLine(18, 76, SidebarWidth-18, 76, rl.Gray)

	// Parameter sliders.
	y := float32(92)

	slider(&simulation.MaxSpeed, 1000, "Max Speed", &y)
	slider(&simulation.MaxForce, 400, "Max Force", &y)
	slider(&simulation.DetectionRange, 200, "Detection Range", &y)
	slider(&simulation.DesiredSeparation, 200, "Desired Separation", &y)
	slider(&simulation.CohesionValue, 1000, "Cohesion", &y)
	slider(&simulation.SeparationValue, 1000, "Separation", &y)
	slider(&simulation.AlignmentValue, 300, "Alignment", &y)

	// Grid optimisation checkbox.
	check := rl.Rectangle{
		X:      left,
		Y:      y + 3,
		Width:  18,
		Height: 18,
	}

	rg.CheckBox(
		check,
		"",
		&simulation.DoGridOptimize,
	)

	rl.DrawText(
		"Use Grid Optimization",
		int32(left+27),
		int32(y+5),
		14,
		rl.RayWhite,
	)

	// show current mode
	modeLabel := "Serial Mode"
	if simulation.IsParallel {
		modeLabel = "Parallel Mode"
	}
	check.Y += 22
	rg.CheckBox(check, "", &simulation.IsParallel)
	rl.DrawText(modeLabel, int32(left+27), int32(check.Y+2), 14, rl.RayWhite)
}

// DrawGUI draws the boid count slider, map information and team credits.
func DrawGUI(setBoidsAmount func(int32)) {
	y := float32(452)

	// Boid count.
	rl.DrawText(
		"Amount of Boids",
		int32(left),
		int32(y),
		14,
		rl.RayWhite,
	)

	rl.DrawText(
		fmt.Sprintf("%d", simulation.BoidsAmount),
		210,
		int32(y),
		14,
		rl.RayWhite,
	)

	temp := float32(simulation.BoidsAmount)

	rg.Slider(
		rl.Rectangle{
			X:      left,
			Y:      y + 23,
			Width:  sliderWidth,
			Height: 16,
		},
		"",
		"",
		&temp,
		0,
		1000,
	)

	if int32(temp) != simulation.BoidsAmount {
		setBoidsAmount(int32(temp))
	}

	// Map switching instructions.
	rl.DrawRectangle(
		18,
		505,
		244,
		32,
		rl.Color{R: 53, G: 70, B: 89, A: 255},
	)

	rl.DrawText(
		"M  -  Switch Map",
		35,
		514,
		14,
		rl.RayWhite,
	)

	// Current map name.
	rl.DrawText(
		simulation.GetMapName(),
		18,
		548,
		12,
		rl.LightGray,
	)

	// Team credits divider.
	rl.DrawLine(
		18,
		580,
		SidebarWidth-18,
		580,
		rl.Gray,
	)

	// Team credits.
	rl.DrawText(
		"PROJECT TEAM",
		18,
		591,
		13,
		rl.LightGray,
	)

	rl.DrawText(
		"Dorian  |  Amelia",
		18,
		615,
		15,
		rl.RayWhite,
	)

	rl.DrawText(
		"Mark    |  Adam",
		18,
		639,
		15,
		rl.RayWhite,
	)

	rl.DrawText(
		"SETU  -  Concurrent Development",
		18,
		675,
		10,
		rl.LightGray,
	)
}

// slider draws a labelled parameter slider.
func slider(value *float32, max float32, label string, y *float32) {
	// Parameter label.
	rl.DrawText(
		label,
		int32(left),
		int32(*y),
		14,
		rl.RayWhite,
	)

	// Current parameter value.
	valueText := fmt.Sprintf("%.0f", *value)

	rl.DrawText(
		valueText,
		SidebarWidth-18-rl.MeasureText(valueText, 14),
		int32(*y),
		14,
		rl.RayWhite,
	)

	// Slider control.
	rg.Slider(
		rl.Rectangle{
			X:      left,
			Y:      *y + 22,
			Width:  sliderWidth,
			Height: 16,
		},
		"",
		"",
		value,
		0,
		max,
	)

	*y += 45
}

// DrawGrid draws grid lines using DetectionRange as the spacing. The spacing must be greater than zero so the loops can finish.
func DrawGrid() {
	if simulation.DetectionRange <= 0 {
		return
	}

	// Vertical grid lines.
	for x := simulation.DetectionRange; x < simulation.Bounds.Width; x += simulation.DetectionRange {
		rl.DrawLine(
			int32(x),
			0,
			int32(x),
			simulation.ScreenY,
			rl.LightGray,
		)
	}

	// Horizontal grid lines.
	for y := simulation.DetectionRange; y < simulation.Bounds.Height; y += simulation.DetectionRange {
		rl.DrawLine(
			0,
			int32(y),
			simulation.ScreenX,
			int32(y),
			rl.LightGray,
		)
	}
}

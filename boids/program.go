package main

import (
	"boids/simulation"
	"boids/ui"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// main opens the window and runs the program.
func main() {
	// open window with chosen size
	rl.InitWindow(simulation.ScreenX+ui.SidebarWidth, simulation.ScreenY, "Boids")

	// close when main ends
	defer rl.CloseWindow()

	//create boids
	spawnBoids()

	// create render texture
	var target rl.RenderTexture2D = rl.LoadRenderTexture(simulation.ScreenX, simulation.ScreenY)

	// keep updating and drawing until window is closed
	for !rl.WindowShouldClose() {
		// Turn off grid optimisation if the detection range is zero.
		if simulation.DetectionRange < 1 {
			simulation.DoGridOptimize = false
		}

		// Press M to switch between free flocking, obstacles and the maze.
		if rl.IsKeyPressed(rl.KeyM) {
			simulation.ToggleMap()
			simulation.RefreshBoidGrid()
		}

		// Get the frame time in seconds.
		dt := rl.GetFrameTime()

		switch simulation.IsParallel {
		case true:
			// concurrent update function
			simulation.UpdateBoids(dt)
		case false:
			// Update each boid one at a time.
			for _, b := range simulation.Boids {
				simulation.BoidAction(b, dt)
			}
		}

		// Start drawing to texture
		rl.BeginTextureMode(target)

		// Clear the screen for the new frame.
		rl.ClearBackground(rl.RayWhite)

		// Show the grid if enabled.
		if simulation.DoGridOptimize {
			ui.DrawGrid()
		}

		// Draw the selected map before drawing the boids.
		simulation.DrawObstacles()

		for i, b := range simulation.Boids {
			// Highlight the first boid and print its velocity.
			if i == 0 {
				simulation.DrawBoid(b, true)
			} else {
				// Draw the rest normally.
				simulation.DrawBoid(b, false)
			}
		}

		// stop drawing to texture
		rl.EndTextureMode()

		// start drawing to screen
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// draw computed texture to screen
		rl.DrawTextureRec(target.Texture, rl.Rectangle{0.0, 0.0, float32(target.Texture.Width), float32(-target.Texture.Height)}, rl.Vector2{float32(ui.SidebarWidth), 0.0}, rl.White)

		// Draw the GUI on the screen so the sidebar can receive mouse input.
		oldDetectionRange := simulation.DetectionRange
		oldGridOptimize := simulation.DoGridOptimize
		ui.DrawSliders()
		ui.DrawGUI(SetBoidsAmount)

		// Update grid cells if the detection range or grid setting changes.
		if simulation.DetectionRange != oldDetectionRange || simulation.DoGridOptimize != oldGridOptimize {
			simulation.RefreshBoidGrid()
		}

		// Show the FPS.
		rl.DrawFPS(ui.SidebarWidth+simulation.Bounds.ToInt32().Width-120, 20)
		rl.EndDrawing()
	}
}

// SetBoidsAmount sets the population of boids for the simulation
func SetBoidsAmount(newAmount int32) {
	simulation.BoidsAmount = newAmount
	spawnBoids()
}

// spawnBoids creates boids with random starting positions.
func spawnBoids() {
	// Clear the old boids.
	simulation.Boids = []*simulation.Boid{}
	// Create the chosen number of boids
	for range simulation.BoidsAmount {
		simulation.Boids = append(simulation.Boids,
			&simulation.Boid{
				// Start in a random spot within 150 pixels right and down from the centre.
				Position: rl.Vector2{X: rand.Float32()*150 + simulation.Bounds.Width*0.5, Y: rand.Float32()*150 + simulation.Bounds.Height*0.5},
			},
		)
	}

	// Move new boids out of any walls and set their grid cells.
	simulation.PrepareBoidsForMap()
}

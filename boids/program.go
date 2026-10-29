package main

import (
	"boids/simulation"
	"boids/ui"
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// main opens the window and runs the program.
func main() {
	// open window with chosen size
	rl.InitWindow(simulation.ScreenX, simulation.ScreenY, "Boids")

	// close when main ends
	defer rl.CloseWindow()

	//create boids
	spawnBoids()

	// create render texture
	var target rl.RenderTexture2D = rl.LoadRenderTexture(simulation.ScreenX, simulation.ScreenY)

	// keep updating and drawing until window is closed
	for !rl.WindowShouldClose() {
		// Turn off grid optimisation if the detection range is zero.
		if simulation.DetectionRange == 0 {
			simulation.DoGridOptimize = false
		}

		// Get the frame time in seconds.
		dt := rl.GetFrameTime()

		// TODO(concurrency, step 1): Replace this loop with simulation.UpdateBoids(dt) in update.go
		// That function should:
		//   1. Copy every *Boid value into a snapshot before any boid moves
		//   2. Calculate each boid's next value from that snapshot, in parallel
		//   3. Wait for all calculations to finish
		//   4. Copy the calculated values back into Boids before drawing below
		// Keep drawing and all raylib calls on the main thread here
		for _, b := range simulation.Boids {
			simulation.BoidAction(b, dt)
		}

		// Start drawing to texture
		rl.BeginTextureMode(target)

		// Clear the screen for the new frame.
		rl.ClearBackground(rl.RayWhite)

		// Show the grid if enabled.
		if simulation.DoGridOptimize {
			ui.DrawGrid()
		}

		for i, b := range simulation.Boids {
			// Highlight the first boid and print its velocity.
			if i == 0 {
				simulation.DrawBoid(b, true)
				fmt.Println(b.Velocity)
			} else {
				// Draw the rest normally.
				simulation.DrawBoid(b, false)
			}
		}

		ui.InitGui()
		ui.DrawGUI(SetBoidsAmount)

		// stop drawing to texture
		rl.EndTextureMode()

		// start drawing to screen
		rl.BeginDrawing()

		// draw computed texture to screen
		rl.DrawTextureRec(target.Texture, rl.Rectangle{0.0, 0.0, float32(target.Texture.Width), float32(-target.Texture.Height)}, rl.Vector2{0.0, 0.0}, rl.White)

		// Show the FPS.
		rl.DrawFPS(simulation.Bounds.ToInt32().Width-80, 20)
		rl.EndDrawing()
	}
}

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
}

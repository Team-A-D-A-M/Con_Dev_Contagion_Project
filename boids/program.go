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
	// Open window with chosen size.
	rl.InitWindow(simulation.ScreenX, simulation.ScreenY, "Boids")

	// Close when main ends.
	defer rl.CloseWindow()

	// Create boids.
	spawnBoids()

	// Create render texture.
	var target rl.RenderTexture2D = rl.LoadRenderTexture(
		simulation.ScreenX,
		simulation.ScreenY,
	)

	defer rl.UnloadRenderTexture(target)

	// Keep updating and drawing until window is closed.
	for !rl.WindowShouldClose() {

		// Press M to switch between maps.
		if rl.IsKeyPressed(rl.KeyM) {
			simulation.ToggleMap()
		}

		// Turn off grid optimisation if detection range is zero.
		if simulation.DetectionRange == 0 {
			simulation.DoGridOptimize = false
		}

		// Get frame time in seconds.
		dt := rl.GetFrameTime()

		// Update each boid.
		for _, b := range simulation.Boids {
			simulation.BoidAction(b, dt)
		}

		// Start drawing to texture.
		rl.BeginTextureMode(target)

		// Clear background.
		rl.ClearBackground(rl.RayWhite)

		// Draw current map title.
		rl.DrawText(
			simulation.GetMapName(),
			20,
			15,
			22,
			rl.DarkGray,
		)

		// Draw simulation border.
		rl.DrawRectangleLines(
			5,
			50,
			simulation.ScreenX-10,
			simulation.ScreenY-55,
			rl.LightGray,
		)

		// Draw grid if enabled.
		if simulation.DoGridOptimize {
			ui.DrawGrid()
		}

		// Draw obstacles belonging to the selected map.
		simulation.DrawObstacles()

		// Draw boids.
		for i, b := range simulation.Boids {

			// Highlight first boid and print its velocity.
			if i == 0 {
				simulation.DrawBoid(b, true)
				fmt.Println(b.Velocity)
			} else {
				simulation.DrawBoid(b, false)
			}
		}

		// Draw GUI.
		ui.InitGui()
		ui.DrawGUI(SetBoidsAmount)

		// Stop drawing to texture.
		rl.EndTextureMode()

		// Start drawing to screen.
		rl.BeginDrawing()

		// Draw computed texture to screen.
		rl.DrawTextureRec(
			target.Texture,
			rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  float32(target.Texture.Width),
				Height: float32(-target.Texture.Height),
			},
			rl.Vector2{X: 0, Y: 0},
			rl.White,
		)

		// Show FPS.
		rl.DrawFPS(
			simulation.Bounds.ToInt32().Width-80,
			20,
		)

		rl.EndDrawing()
	}
}

// SetBoidsAmount changes the number of boids.
func SetBoidsAmount(newAmount int32) {
	simulation.BoidsAmount = newAmount
	spawnBoids()
}

// spawnBoids creates boids with random starting positions.
func spawnBoids() {

	// Clear old boids.
	simulation.Boids = []*simulation.Boid{}

	// Create the chosen number of boids.
	for range simulation.BoidsAmount {
		simulation.Boids = append(
			simulation.Boids,
			&simulation.Boid{
				// Start in a random position around the centre.
				Position: rl.Vector2{
					X: rand.Float32()*150 + simulation.Bounds.Width*0.5,
					Y: rand.Float32()*150 + simulation.Bounds.Height*0.5,
				},
			},
		)
	}
}

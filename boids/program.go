package main

import (
	"boids/simulation"
	"boids/ui"
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(simulation.ScreenX, simulation.ScreenY, "Boids")
	defer rl.CloseWindow()

	spawnBoids()
	for !rl.WindowShouldClose() {
		if simulation.DetectionRange == 0 {
			simulation.DoGridOptimize = false
		}
		for _, b := range simulation.Boids {
			simulation.BoidAction(b)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		if simulation.DoGridOptimize {
			ui.DrawGrid()
		}

		for i, b := range simulation.Boids {
			if i == 0 {
				simulation.DrawBoid(b, true)
				fmt.Println(b.Velocity)
			} else {
				simulation.DrawBoid(b, false)
			}
		}

		ui.InitGui()

		rl.DrawFPS(simulation.Bounds.ToInt32().Width-80, 20)
		rl.EndDrawing()
	}
}

func SetBoidsAmount(newAmount int32) {
	simulation.BoidsAmount = newAmount
	spawnBoids()
}

func spawnBoids() {
	simulation.Boids = []*simulation.Boid{}
	for range simulation.BoidsAmount {
		simulation.Boids = append(simulation.Boids,
			&simulation.Boid{
				Position: rl.Vector2{X: rand.Float32()*150 + simulation.Bounds.Width*0.5, Y: rand.Float32()*150 + simulation.Bounds.Height*0.5},
			},
		)
	}
}

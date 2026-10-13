package main

import (
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// globals

var boidsAmount int32 = 400
var doGridOptimize bool = true

var boids = []*Boid{}

func main() {
	rl.InitWindow(1100, 700, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	//rl.SetTargetFPS(60)

	spawnBoids()

	for !rl.WindowShouldClose() {
		if detectionRange == 0 {
			doGridOptimize = false
		}
		for _, b := range boids {
			boidAction(b)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		if doGridOptimize {
			drawGrid()
		}

		for i, b := range boids {
			if i == 0 {
				drawBoid(b, true)
				fmt.Println(b.velocity)
			} else {
				drawBoid(b, false)
			}
		}

		initGui()

		rl.DrawFPS(bounds.ToInt32().Width-80, 20)
		rl.EndDrawing()
	}
}

func setBoidsAmount(newAmount int32) {
	boidsAmount = newAmount
	spawnBoids()
}

func spawnBoids() {
	boids = []*Boid{}
	for range boidsAmount {
		boids = append(boids,
			&Boid{
				position: rl.Vector2{X: rand.Float32()*150 + bounds.Width*0.5, Y: rand.Float32()*150 + bounds.Height*0.5},
			},
		)
	}
}

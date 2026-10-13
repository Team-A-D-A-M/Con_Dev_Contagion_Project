package main

import (
	"fmt"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// globals

var boidsAmount int32 = 400

var boids = []*Boid{}

func main() {
	fmt.Println(rl.Vector2Angle(up2d, rl.Vector2{X: 1.0, Y: 0.0}))
	rl.InitWindow(1100, 700, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for range boidsAmount {
		boids = append(boids,
			&Boid{
				position: rl.Vector2{X: rand.Float32()*100 + 200, Y: rand.Float32()*100 + 200},
			},
		)
	}

	for !rl.WindowShouldClose() {
		for _, b := range boids {
			boidAction(b)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for i, b := range boids {
			if i == 0 {
				drawBoid(b, true)
				fmt.Println(b.velocity)
			} else {
				drawBoid(b, false)
			}
		}
		initGui()
		rl.EndDrawing()
	}
}

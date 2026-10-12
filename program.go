package condevcontagionproject

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// globals
var bounds = rl.Rectangle{
	X:      300,
	Y:      100,
	Width:  500,
	Height: 500,
}

func main() {
	rl.InitWindow(1100, 700, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawRectangleLinesEx(bounds, 2, rl.Blue)

		rl.EndDrawing()
	}
}

func moveAllAgents() {

}

func rollInfect() {

}

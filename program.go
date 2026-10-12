package main

import (
	rg "github.com/gen2brain/raylib-go/raygui"
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

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rg.Button(bounds, "text")

		rl.ClearBackground(rl.RayWhite)

		rl.EndDrawing()
	}
}

func moveAllAgents() {

}

func rollInfect() {

}

package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var bounds rl.Rectangle = rl.Rectangle{X: 0, Y: 0, Width: 1100, Height: 700}

var up2d rl.Vector2 = rl.Vector2{X: 0.0, Y: -1.0}
var down2d rl.Vector2 = rl.Vector2{X: 0.0, Y: 1.0}
var left2d rl.Vector2 = rl.Vector2{X: -1.0, Y: 0.0}
var right2d rl.Vector2 = rl.Vector2{X: 1.0, Y: 0.0}

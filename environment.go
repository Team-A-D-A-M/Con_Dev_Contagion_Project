package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// screen boundaries
var bounds rl.Rectangle = rl.Rectangle{X: 0, Y: 0, Width: float32(screenX), Height: float32(screenY)}

// spatial directions
var up2d rl.Vector2 = rl.Vector2{X: 0.0, Y: -1.0}
var down2d rl.Vector2 = rl.Vector2{X: 0.0, Y: 1.0}
var left2d rl.Vector2 = rl.Vector2{X: -1.0, Y: 0.0}
var right2d rl.Vector2 = rl.Vector2{X: 1.0, Y: 0.0}

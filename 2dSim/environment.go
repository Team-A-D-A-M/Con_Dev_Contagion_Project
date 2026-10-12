package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var boids = []*Boid{}
var bounds rl.Rectangle = rl.Rectangle{X: 0, Y: 0, Width: 400, Height: 400}

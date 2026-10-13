package main

// boid parameters
var detectionRange float32 = 80.0
var maxSpeed float32 = 600.0
var maxForce float32 = 200.0
var desiredSeparation float32 = 24.0

// boid rules
var cohesionValue float32 = 100
var separationValue float32 = 250
var alignmentValue float32 = 150

// simulation parameters
var boidsAmount int32 = 400
var doGridOptimize bool = true
var screenX int32 = 1100
var screenY int32 = 700

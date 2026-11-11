package simulation

// boid parameters
var DetectionRange float32 = 80.0
var MaxSpeed float32 = 600.0
var MaxForce float32 = 200.0
var DesiredSeparation float32 = 24.0

// boid rules
var CohesionValue float32 = 100
var SeparationValue float32 = 250
var AlignmentValue float32 = 150

// simulation parameters
var BoidsAmount int32 = 400
var DoGridOptimize bool = true
var ScreenX int32 = 1100
var ScreenY int32 = 700

// active simulation boid population
var Boids []*Boid

// number of worker pools to spawn
var Workers = 10

var IsParallel = true

var ParallelOrSerial = "Parallel"

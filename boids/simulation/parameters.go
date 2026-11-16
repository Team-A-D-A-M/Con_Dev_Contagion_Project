package simulation

// DetectionRange the maximum radius in pixels a boid can detect another boid.
var DetectionRange float32 = 80.0

// MaxSpeed the maximum velocity of a boid.
var MaxSpeed float32 = 600.0

// MaxForce the maximum steering force for a boid.
var MaxForce float32 = 200.0

// DesiredSeparation the minimum personal space between any boid.
var DesiredSeparation float32 = 24.0

// CohesionValue sets the strength of the force that attracts boids together.
var CohesionValue float32 = 100

// SeparationValue sets the strength of the force that repels boids from each other.
var SeparationValue float32 = 250

// AlignmentValue sets how strongly a boid tries to match the direction of a flock.
var AlignmentValue float32 = 150

// simulation parameters
// BoidsAmount the number of boids per simulation
var BoidsAmount int32 = 400

// DoGridOptimize toggles grid optimisation to speed up neighbour look-up times
var DoGridOptimize bool = true

// ScreenX the length of the screen
var ScreenX int32 = 1100

// ScreenY the height of the screen
var ScreenY int32 = 700

// Boids active simulation boid population
var Boids []*Boid

// Workers number of worker pools to spawn
var Workers = 10

// IsParallel boolean val sets parallel implementation as default
var IsParallel = true

package simulation

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// applyObstacleSteering uses obstacle avoidance when there is a wall ahead.
func applyObstacleSteering(b *Boid, steering rl.Vector2) rl.Vector2 {
	avoidance := obstacleAvoidance(b)

	// Use a stronger steering force so the boid can turn before hitting the obstacle.
	if avoidance.Length() > 0 {
		return limitVector(avoidance, MaxForce*7)
	}

	// Keep the original flocking steering when there is no obstacle ahead.
	return steering
}

// RefreshBoidGrid updates the grid cells after changing the map or grid settings.
func RefreshBoidGrid() {
	// Do not divide by the detection range when grid optimisation is disabled.
	if !DoGridOptimize || DetectionRange < 1 {
		return
	}

	for _, b := range Boids {
		updateBoidGrid(b)
	}
}

// PrepareBoidsForMap moves new boids out of walls and sets their grid cells.
func PrepareBoidsForMap() {
	for _, b := range Boids {
		ResolveObstacleCollision(b, b.Position)
	}
	RefreshBoidGrid()
}

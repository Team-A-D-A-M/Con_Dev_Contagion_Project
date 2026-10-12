package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var detectionRange float32 = 50.0
var maxSpeed float32 = 100.0
var maxForce float32 = 50.0
var desiredSeparation float32 = 20.0

var cohesionValue float32 = 1.0
var separationValue float32 = 1.0
var alignmentValue float32 = 1.0

type Boid struct {
	position rl.Vector2
	velocity rl.Vector2
}

type NearbyBoid struct {
	boid     *Boid
	distance float32
}

// Steer towards the average position of nearby boids.
func boidCohesion(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	center := rl.Vector2Zero()

	if len(nearby) == 0 {
		return rl.Vector2Zero()
	}

	for _, n := range nearby {
		center = rl.Vector2Add(center, n.boid.position)
	}

	center = rl.Vector2Scale(center, 1.0/float32(len(nearby)))
	desiredVelocity := rl.Vector2Subtract(center, b.position)

	if desiredVelocity.Length() > 0 {
		desiredVelocity = rl.Vector2Scale(
			rl.Vector2Normalize(desiredVelocity),
			maxSpeed,
		)
	}

	return rl.Vector2Subtract(desiredVelocity, b.velocity)
}

// Steer away from nearby boids that are too close.
func boidSeparation(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	steer := rl.Vector2Zero()
	count := 0

	for _, n := range nearby {
		distance := n.distance

		// Coincident boids have no well-defined separation direction.
		if distance <= 0.001 || distance >= desiredSeparation {
			continue
		}

		difference := rl.Vector2Subtract(b.position, n.boid.position)
		difference = rl.Vector2Scale(
			rl.Vector2Normalize(difference),
			1.0/distance,
		)

		steer = rl.Vector2Add(steer, difference)
		count++
	}

	if count == 0 {
		return rl.Vector2Zero()
	}

	steer = rl.Vector2Scale(steer, 1.0/float32(count))

	if steer.Length() > 0 {
		desiredVelocity := rl.Vector2Scale(
			rl.Vector2Normalize(steer),
			maxSpeed,
		)
		return rl.Vector2Subtract(desiredVelocity, b.velocity)
	}

	return rl.Vector2Zero()
}

// Steer towards the average velocity of nearby boids.
func boidAlignment(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	if len(nearby) == 0 {
		return rl.Vector2Zero()
	}

	avgVelocity := rl.Vector2Zero()

	for _, n := range nearby {
		avgVelocity = rl.Vector2Add(avgVelocity, n.boid.velocity)
	}

	avgVelocity = rl.Vector2Scale(
		avgVelocity,
		1.0/float32(len(nearby)),
	)

	if avgVelocity.Length() == 0 {
		return rl.Vector2Zero()
	}

	desiredVelocity := rl.Vector2Scale(
		rl.Vector2Normalize(avgVelocity),
		maxSpeed,
	)

	return rl.Vector2Subtract(desiredVelocity, b.velocity)
}

// Limit a vector's magnitude.
func limitVector(v rl.Vector2, max float32) rl.Vector2 {
	length := v.Length()

	if length > max {
		return rl.Vector2Scale(rl.Vector2Normalize(v), max)
	}

	return v
}

func boidAction(b *Boid) {
	nearby := getNearbyBoids(b)

	cohesion := boidCohesion(b, nearby)
	separation := boidSeparation(b, nearby)
	alignment := boidAlignment(b, nearby)

	steering := rl.Vector2Add(
		rl.Vector2Scale(cohesion, cohesionValue),
		rl.Vector2Add(
			rl.Vector2Scale(separation, separationValue),
			rl.Vector2Scale(alignment, alignmentValue),
		),
	)

	steering = limitVector(steering, maxForce)

	dt := rl.GetFrameTime()

	b.velocity = rl.Vector2Add(
		b.velocity,
		rl.Vector2Scale(steering, dt),
	)
	b.velocity = limitVector(b.velocity, maxSpeed)

	b.position = rl.Vector2Add(
		b.position,
		rl.Vector2Scale(b.velocity, dt),
	)
}

func getNearbyBoids(b *Boid) []*NearbyBoid {
	nearbyBoids := []*NearbyBoid{}

	for _, other := range boids {
		if other == b {
			continue
		}

		distance := rl.Vector2Distance(b.position, other.position)

		if distance < detectionRange {
			nearbyBoids = append(nearbyBoids, &NearbyBoid{
				boid:     other,
				distance: distance,
			})
		}
	}

	return nearbyBoids
}

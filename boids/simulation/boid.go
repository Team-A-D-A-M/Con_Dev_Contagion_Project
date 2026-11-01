package simulation

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// GridCell stores a grid's location
type GridCell struct {
	X, Y int
}

// Boid stores a boid's position, velocity and location
type Boid struct {
	Position rl.Vector2
	Velocity rl.Vector2
	GridCell GridCell
}

// NearbyBoid stores a nearby boid and it's distance
type NearbyBoid struct {
	Boid     *Boid
	Distance float32
}

// Steer towards the average position of nearby boids.
func boidCohesion(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	center := rl.Vector2Zero()

	//return no steering when there are no nearby boids
	if len(nearby) == 0 {
		return rl.Vector2Zero()
	}

	// Add positions of nearby boids
	for _, n := range nearby {
		center = rl.Vector2Add(center, n.Boid.Position)
	}

	// calculate average position of nearby boids
	center = rl.Vector2Scale(center, 1.0/float32(len(nearby)))
	// Find the direction towards the centre of nearby boids.
	desiredVelocity := rl.Vector2Subtract(center, b.Position)

	// Set the desired speed to MaxSpeed if the direction is not zero.
	if desiredVelocity.Length() > 0 {
		// Normalize makes the vector length 1 without changing its direction.
		// Scale multiplies it by MaxSpeed to set the desired speed.
		desiredVelocity = rl.Vector2Scale(
			rl.Vector2Normalize(desiredVelocity),
			MaxSpeed,
		)
	}

	// Calculate steering from the current velocity.
	return rl.Vector2Subtract(desiredVelocity, b.Velocity)
}

// Steer away from nearby boids that are too close.
func boidSeparation(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	steer := rl.Vector2Zero()
	count := 0

	// Check each nearby boid.
	for _, n := range nearby {
		distance := n.Distance

		// Skip overlapping boids and boids that are far enough away
		if distance <= 0.001 || distance >= DesiredSeparation {
			continue
		}

		// Find the direction away from the nearby boid.
		difference := rl.Vector2Subtract(b.Position, n.Boid.Position)

		// Give closer boids a stronger influence.
		difference = rl.Vector2Scale(
			rl.Vector2Normalize(difference),
			1.0/distance,
		)

		// Add the vector pointing away from this neighbour to the total.
		steer = rl.Vector2Add(steer, difference)
		count++
	}

	// Return zero steering if no neighbours were included.
	if count == 0 {
		return rl.Vector2Zero()
	}

	// Divide the total by the number of included neighbours.
	steer = rl.Vector2Scale(steer, 1.0/float32(count))

	// Check that there is a direction to steer towards.
	if steer.Length() > 0 {
		// Set the desired velocity in this direction at MaxSpeed.
		desiredVelocity := rl.Vector2Scale(
			rl.Vector2Normalize(steer),
			MaxSpeed,
		)
		// Subtract the current velocity to calculate the steering needed.
		return rl.Vector2Subtract(desiredVelocity, b.Velocity)
	}

	// Return zero steering if the vectors cancel each other out.
	return rl.Vector2Zero()
}

// Steer towards the average velocity of nearby boids.
func boidAlignment(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	if len(nearby) == 0 {
		return rl.Vector2Zero()
	} //no steering if no nearby boids

	avgVelocity := rl.Vector2Zero()

	// Add up the nearby boid's velocities
	for _, n := range nearby {
		avgVelocity = rl.Vector2Add(avgVelocity, n.Boid.Velocity)
	}

	// Divide the total by the number of nearby boids.
	avgVelocity = rl.Vector2Scale(
		avgVelocity,
		1.0/float32(len(nearby)),
	)

	// Stop here if the average vector is zero.
	if avgVelocity.Length() == 0 {
		return rl.Vector2Zero()
	}

	// Use the average direction with a target speed of MaxSpeed.
	desiredVelocity := rl.Vector2Scale(
		rl.Vector2Normalize(avgVelocity),
		MaxSpeed,
	)

	// Subtract the current velocity to calculate the steering needed.
	return rl.Vector2Subtract(desiredVelocity, b.Velocity)
}

// BoidAction updates the boid's movement using the 3 flocking rules
func BoidActionConcurrent(index int, snapshot []Boid, dt float32) Boid {
	b := snapshot[index]

	// concurrent nearby boids method so it doesnt read global boids
	nearby := getNearbyBoidsConcurrent(index, snapshot)

	// apply the three flocking rules
	cohesion := boidCohesion(&b, nearby)
	separation := boidSeparation(&b, nearby)
	alignment := boidAlignment(&b, nearby)

	boidMovement(&b, dt, cohesion, separation, alignment)

	// Update the grid cell for the new position.
	if DoGridOptimize {
		updateBoidGrid(&b)
	}

	return b
}

func BoidAction(b *Boid, dt float32) {
	// Find neighbours within the detection range.
	nearby := getNearbyBoids(b)

	// Calculate steering towards the group's centre.
	cohesion := boidCohesion(b, nearby)

	// Calculate steering away from boids that are too close.
	separation := boidSeparation(b, nearby)

	// Calculate steering in the neighbours' average direction.
	alignment := boidAlignment(b, nearby)

	boidMovement(b, dt, cohesion, separation, alignment)

	// Update which grid cell the boid is in.
	if DoGridOptimize {
		updateBoidGrid(b)
	}
}

// function with extracted common code for boid action used by parallel and serial versions
func boidMovement(b *Boid, dt float32, cohesion rl.Vector2, separation rl.Vector2, alignment rl.Vector2) {
	// Multiply each rule by its weight and add the results.
	steering := rl.Vector2Add(
		rl.Vector2Scale(cohesion, CohesionValue),
		rl.Vector2Add(
			rl.Vector2Scale(separation, SeparationValue),
			rl.Vector2Scale(alignment, AlignmentValue),
		),
	)

	// Keep the steering force within MaxForce.
	steering = limitVector(steering, MaxForce)

	// Change the velocity based on steering and frame time.
	b.Velocity = rl.Vector2Add(
		b.Velocity,
		rl.Vector2Scale(steering, dt),
	)
	// Stop the boid from moving faster than MaxSpeed.
	b.Velocity = limitVector(b.Velocity, MaxSpeed)

	// Move the boid based on its velocity and frame time.
	b.Position = rl.Vector2Add(
		b.Position,
		rl.Vector2Scale(b.Velocity, dt),
	)

	// Move boids crossing the left or right edge to the opposite side.
	if b.Position.X < 0.0 {
		b.Position.X += Bounds.Width
	} else if b.Position.X > Bounds.Width {
		b.Position.X -= Bounds.Width
	}
	// Move boids crossing the top or bottom edge to the opposite side.
	if b.Position.Y < 0.0 {
		b.Position.Y += Bounds.Height
	} else if b.Position.Y > Bounds.Height {
		b.Position.Y -= Bounds.Height
	}
}

// updateBoidGrid updates the grid cell containing the boid.
func updateBoidGrid(b *Boid) {
	// Divide each position coordinate by the cell size and round down.
	b.GridCell.X = floorDiv(b.Position.X, DetectionRange)
	b.GridCell.Y = floorDiv(b.Position.Y, DetectionRange)
}

func getNearbyBoidsConcurrent(index int, snapshot []Boid) []*NearbyBoid {
	nearbyBoids := make([]*NearbyBoid, 0)

	for i := range snapshot {

		if i == index {
			continue
		}

		other := &snapshot[i]

		// Find how many grid cells apart the boids are.
		diffX := snapshot[index].GridCell.X - other.GridCell.X
		diffY := snapshot[index].GridCell.Y - other.GridCell.Y
		// When optimising, only check the same cell and adjacent cells.
		if (diffX < -1 || diffX > 1 || diffY < -1 || diffY > 1) && DoGridOptimize {
			continue
		}

		// Measure the distance between the boids.
		distance := rl.Vector2Distance(snapshot[index].Position, other.Position)

		// Add the neighbour and its distance if it is within range.
		if distance < DetectionRange {
			nearbyBoids = append(nearbyBoids, &NearbyBoid{
				Boid:     other,
				Distance: distance,
			})
		}
	}

	return nearbyBoids
}

// getNearbyBoids finds other boids within the detection range.
func getNearbyBoids(b *Boid) []*NearbyBoid {
	// Start with an empty list.
	nearbyBoids := []*NearbyBoid{}

	for _, other := range Boids {
		// Do not count the boid itself as a neighbour.
		if other == b {
			continue
		}

		// Find how many grid cells apart the boids are.
		diffX := b.GridCell.X - other.GridCell.X
		diffY := b.GridCell.Y - other.GridCell.Y
		// When optimising, only check the same cell and adjacent cells.
		if (diffX < -1 || diffX > 1 || diffY < -1 || diffY > 1) && DoGridOptimize {
			continue
		}

		// Measure the distance between the boids.
		distance := rl.Vector2Distance(b.Position, other.Position)

		// Add the neighbour and its distance if it is within range.
		if distance < DetectionRange {
			nearbyBoids = append(nearbyBoids, &NearbyBoid{
				Boid:     other,
				Distance: distance,
			})
		}
	}

	// Return the neighbours found.
	return nearbyBoids
}

// lineLength sets how long the direction line is.
var lineLength float32 = 24.0

// DrawBoid draws the boid and highlights it when requested.
func DrawBoid(b *Boid, highlight bool) {
	// Use blue unless the boid is highlighted.
	col := rl.Blue
	if highlight {
		col = rl.Red
	}

	vel := rl.Vector2Normalize(b.Velocity)
	// Show the highlighted boid's detection range as a circle.
	if highlight {
		rl.DrawCircleLines(int32(b.Position.X), int32(b.Position.Y), DetectionRange, rl.LightGray)
	}
	// Draw a line showing which way the boid is moving.
	rl.DrawLine(int32(b.Position.X), int32(b.Position.Y), int32(b.Position.X+vel.X*lineLength), int32(b.Position.Y+vel.Y*lineLength), rl.LightGray)
	// Draw the boid as a triangle facing its movement direction.
	rl.DrawPoly(
		b.Position,
		3,
		6,
		rl.Vector2Angle(right2d, b.Velocity)*rl.Rad2deg,
		col,
	)
}

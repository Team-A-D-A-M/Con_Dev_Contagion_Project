package simulation

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// GridCell stores a grid's location.
type GridCell struct {
	X, Y int
}

// Boid stores a boid's position, velocity and location.
type Boid struct {
	Position rl.Vector2 // position of the boid on screen
	Velocity rl.Vector2 // the boid's movement every frame
	GridCell GridCell   // Which cell in the optimization grid does the boid occupy
}

// NearbyBoid stores a nearby boid and it's distance.
type NearbyBoid struct {
	Boid     *Boid   // the nearby boid that is referenced
	Distance float32 // the boid's distance from the relevant boid
}

// BoidCohesion returns a value steering a boid towards the average position of nearby boids.
func BoidCohesion(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
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

// BoidSeparation returns a value steering a boid away from nearby boids that are too close.
func BoidSeparation(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
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

// BoidAlignment returns a value steering a boid towards the average velocity of nearby boids.
func BoidAlignment(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	// no steering if no nearby boids
	if len(nearby) == 0 {
		return rl.Vector2Zero()
	}

	// Add up the nearby boid's velocities
	avgVelocity := rl.Vector2Zero()
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

// BoidActionConcurrent updates the boid's movement using the 3 flocking rules in parallel.
func BoidActionConcurrent(index int, snapshot []Boid, dt float32) Boid {
	b := snapshot[index]

	// concurrent nearby boids method so it doesnt read global boids
	nearby := getNearbyBoidsConcurrent(index, snapshot)

	// apply the three flocking rules
	cohesion := BoidCohesion(&b, nearby)
	separation := BoidSeparation(&b, nearby)
	alignment := BoidAlignment(&b, nearby)

	BoidMovement(&b, dt, cohesion, separation, alignment)

	return b
}

// BoidAction updates the boid's movement using the 3 flocking rules.
func BoidAction(b *Boid, dt float32) {
	// Find neighbours within the detection range.
	nearby := GetNearbyBoids(b)

	// apply the three flocking rules
	cohesion := BoidCohesion(b, nearby)
	separation := BoidSeparation(b, nearby)
	alignment := BoidAlignment(b, nearby)

	BoidMovement(b, dt, cohesion, separation, alignment)
}

// BoidMovemement applies a per-frame movement based on the boid's 3 rules results.
func BoidMovement(b *Boid, dt float32, cohesion rl.Vector2, separation rl.Vector2, alignment rl.Vector2) {

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

	// Give obstacle avoidance priority when the boid is heading towards a wall.
	steering = applyObstacleSteering(b, steering)

	// Save the position before moving so collisions can be checked along the path.
	previous := b.Position

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

	// Stop the boid from moving through an obstacle before wrapping the position.
	ResolveObstacleCollision(b, previous)

	b.Position = positionWrap(b.Position)

	// Update which grid cell the boid is in.
	if DoGridOptimize {
		updateBoidGrid(b)
	}
}

// UpdateBoidGrid updates the grid cell containing the boid.
func updateBoidGrid(b *Boid) {
	// Divide each position coordinate by the cell size and round down.
	b.GridCell.X = floorDiv(b.Position.X, DetectionRange)
	b.GridCell.Y = floorDiv(b.Position.Y, DetectionRange)
}

// Gets nearby boids within detection range for concurrent implementation.
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

// GetNearbyBoids returns an array of other boids within the detection range of a boid.
func GetNearbyBoids(b *Boid) []*NearbyBoid {
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

// DrawBoid draws the boid and highlights it if set.
func DrawBoid(b *Boid, highlight bool) {

	// set the length of the direction line
	var lineLength float32 = 24.0

	// Use blue unless the boid is highlighted, then use red
	col := rl.Blue
	if highlight {
		col = rl.Red
	}

	// get boid velocity for drawing the visual shape
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

// PositionWrap wraps a Vector2 to the opposite end of the screen if it goes out-of-bounds
func positionWrap(position rl.Vector2) rl.Vector2 {
	// horizontal wrap
	if position.X < 0.0 {
		position.X += Bounds.Width
	} else if position.X > Bounds.Width {
		position.X -= Bounds.Width
	}

	// vertical wrap
	if position.Y < 0.0 {
		position.Y += Bounds.Height
	} else if position.Y > Bounds.Height {
		position.Y -= Bounds.Height
	}

	return position
}

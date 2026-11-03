package simulation

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// GridCell stores a grid's location.
type GridCell struct {
	X, Y int
}

// Boid stores a boid's position, velocity and location.
type Boid struct {
	Position rl.Vector2
	Velocity rl.Vector2
	GridCell GridCell
}

// NearbyBoid stores a nearby boid and its distance.
type NearbyBoid struct {
	Boid     *Boid
	Distance float32
}

// Steer towards the average position of nearby boids.
func boidCohesion(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	center := rl.Vector2Zero()

	if len(nearby) == 0 {
		return rl.Vector2Zero()
	}

	for _, n := range nearby {
		center = rl.Vector2Add(center, n.Boid.Position)
	}

	center = rl.Vector2Scale(center, 1.0/float32(len(nearby)))
	desiredVelocity := rl.Vector2Subtract(center, b.Position)

	if desiredVelocity.Length() > 0 {
		desiredVelocity = rl.Vector2Scale(
			rl.Vector2Normalize(desiredVelocity),
			MaxSpeed,
		)
	}

	return rl.Vector2Subtract(desiredVelocity, b.Velocity)
}

// Steer away from nearby boids that are too close.
func boidSeparation(b *Boid, nearby []*NearbyBoid) rl.Vector2 {
	steer := rl.Vector2Zero()
	count := 0

	for _, n := range nearby {
		distance := n.Distance

		if distance <= 0.001 || distance >= DesiredSeparation {
			continue
		}

		difference := rl.Vector2Subtract(b.Position, n.Boid.Position)

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
			MaxSpeed,
		)

		return rl.Vector2Subtract(desiredVelocity, b.Velocity)
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
		avgVelocity = rl.Vector2Add(avgVelocity, n.Boid.Velocity)
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
		MaxSpeed,
	)

	return rl.Vector2Subtract(desiredVelocity, b.Velocity)
}

// BoidAction updates the boid's movement using the flocking rules.
func BoidAction(b *Boid, dt float32) {
	nearby := getNearbyBoids(b)

	cohesion := boidCohesion(b, nearby)
	separation := boidSeparation(b, nearby)
	alignment := boidAlignment(b, nearby)
	avoidance := obstacleAvoidance(b)

	// Combine the flocking rules.
	steering := rl.Vector2Add(
		rl.Vector2Scale(cohesion, CohesionValue),
		rl.Vector2Add(
			rl.Vector2Scale(separation, SeparationValue),
			rl.Vector2Scale(alignment, AlignmentValue),
		),
	)

	// Give obstacle avoidance priority.
	if avoidance.Length() > 0 {
		steering = rl.Vector2Add(
			steering,
			rl.Vector2Scale(avoidance, 1000),
		)
	}

	steering = limitVector(steering, MaxForce)

	// Update velocity.
	b.Velocity = rl.Vector2Add(
		b.Velocity,
		rl.Vector2Scale(steering, dt),
	)

	b.Velocity = limitVector(b.Velocity, MaxSpeed)

	// Move with collision detection.
	moveBoidWithCollision(b, dt)

	// Wrap around screen edges.
	if b.Position.X < 0 {
		b.Position.X += Bounds.Width
	} else if b.Position.X > Bounds.Width {
		b.Position.X -= Bounds.Width
	}

	if b.Position.Y < 0 {
		b.Position.Y += Bounds.Height
	} else if b.Position.Y > Bounds.Height {
		b.Position.Y -= Bounds.Height
	}

	// Update grid position.
	if DoGridOptimize {
		updateBoidGrid(b)
	}
}

// Move the boid in small steps to prevent it passing through obstacles.
func moveBoidWithCollision(b *Boid, dt float32) {
	distance := b.Velocity.Length() * dt

	// Smaller steps make collision detection more reliable at high speeds.
	steps := int(math.Ceil(float64(distance / 3)))

	if steps < 1 {
		steps = 1
	}

	stepTime := dt / float32(steps)

	for i := 0; i < steps; i++ {
		nextPosition := rl.Vector2Add(
			b.Position,
			rl.Vector2Scale(b.Velocity, stepTime),
		)

		if CurrentMap == 1 {
			nextPosition = resolveMap1Collisions(nextPosition, &b.Velocity)
		} else {
			nextPosition = resolveMazeCollisions(nextPosition, &b.Velocity)
		}

		b.Position = nextPosition
	}
}

// Handle collisions with Map 1 circles and rectangles.
func resolveMap1Collisions(position rl.Vector2, velocity *rl.Vector2) rl.Vector2 {
	const boidRadius float32 = 7

	// Check circular obstacles.
	for _, obstacle := range Obstacles {
		difference := rl.Vector2Subtract(position, obstacle.Position)
		distance := difference.Length()
		minimumDistance := obstacle.Radius + boidRadius

		if distance < minimumDistance {
			var normal rl.Vector2

			if distance > 0.001 {
				normal = rl.Vector2Scale(difference, 1.0/distance)
			} else {
				normal = rl.Vector2{X: 1, Y: 0}
			}

			// Push the boid outside the circle.
			position = rl.Vector2Add(
				obstacle.Position,
				rl.Vector2Scale(normal, minimumDistance),
			)

			// Remove velocity directed into the obstacle.
			inwardSpeed := rl.Vector2DotProduct(*velocity, normal)

			if inwardSpeed < 0 {
				*velocity = rl.Vector2Subtract(
					*velocity,
					rl.Vector2Scale(normal, inwardSpeed),
				)
			}
		}
	}

	// Check rectangular obstacles.
	position = resolveRectangles(position, velocity, RectObstacles)

	return position
}

// Handle collisions with the maze walls.
func resolveMazeCollisions(position rl.Vector2, velocity *rl.Vector2) rl.Vector2 {
	return resolveRectangles(position, velocity, MazeObstacles)
}

// Resolve collisions with a collection of rectangles.
func resolveRectangles(position rl.Vector2, velocity *rl.Vector2, obstacles []RectObstacle) rl.Vector2 {
	const boidRadius float32 = 7

	for _, obstacle := range obstacles {
		r := obstacle.Bounds

		minX := r.X - boidRadius
		maxX := r.X + r.Width + boidRadius
		minY := r.Y - boidRadius
		maxY := r.Y + r.Height + boidRadius

		if position.X < minX ||
			position.X > maxX ||
			position.Y < minY ||
			position.Y > maxY {
			continue
		}

		left := position.X - minX
		right := maxX - position.X
		top := position.Y - minY
		bottom := maxY - position.Y

		smallest := left
		side := 0

		if right < smallest {
			smallest = right
			side = 1
		}

		if top < smallest {
			smallest = top
			side = 2
		}

		if bottom < smallest {
			side = 3
		}

		var normal rl.Vector2

		switch side {
		case 0:
			position.X = minX
			normal = rl.Vector2{X: -1, Y: 0}

		case 1:
			position.X = maxX
			normal = rl.Vector2{X: 1, Y: 0}

		case 2:
			position.Y = minY
			normal = rl.Vector2{X: 0, Y: -1}

		case 3:
			position.Y = maxY
			normal = rl.Vector2{X: 0, Y: 1}
		}

		// Remove velocity directed into the wall.
		inwardSpeed := rl.Vector2DotProduct(*velocity, normal)

		if inwardSpeed < 0 {
			*velocity = rl.Vector2Subtract(
				*velocity,
				rl.Vector2Scale(normal, inwardSpeed),
			)
		}
	}

	return position
}

// updateBoidGrid updates the grid cell containing the boid.
func updateBoidGrid(b *Boid) {
	b.GridCell.X = floorDiv(b.Position.X, DetectionRange)
	b.GridCell.Y = floorDiv(b.Position.Y, DetectionRange)
}

// getNearbyBoids finds other boids within the detection range.
func getNearbyBoids(b *Boid) []*NearbyBoid {
	nearbyBoids := []*NearbyBoid{}

	for _, other := range Boids {
		if other == b {
			continue
		}

		diffX := b.GridCell.X - other.GridCell.X
		diffY := b.GridCell.Y - other.GridCell.Y

		if (diffX < -1 || diffX > 1 || diffY < -1 || diffY > 1) && DoGridOptimize {
			continue
		}

		distance := rl.Vector2Distance(b.Position, other.Position)

		if distance < DetectionRange {
			nearbyBoids = append(nearbyBoids, &NearbyBoid{
				Boid:     other,
				Distance: distance,
			})
		}
	}

	return nearbyBoids
}

// Line length for the boid direction indicator.
var lineLength float32 = 24.0

// DrawBoid draws the boid and highlights it when requested.
func DrawBoid(b *Boid, highlight bool) {
	col := rl.Blue

	if highlight {
		col = rl.Red
	}

	vel := rl.Vector2Normalize(b.Velocity)

	if highlight {
		rl.DrawCircleLines(
			int32(b.Position.X),
			int32(b.Position.Y),
			DetectionRange,
			rl.LightGray,
		)
	}

	rl.DrawLine(
		int32(b.Position.X),
		int32(b.Position.Y),
		int32(b.Position.X+vel.X*lineLength),
		int32(b.Position.Y+vel.Y*lineLength),
		rl.LightGray,
	)

	rl.DrawPoly(
		b.Position,
		3,
		6,
		rl.Vector2Angle(right2d, b.Velocity)*rl.Rad2deg,
		col,
	)
}

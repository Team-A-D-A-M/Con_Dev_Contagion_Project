package simulation

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GridCell struct {
	X, Y int
}

type Boid struct {
	Position rl.Vector2
	Velocity rl.Vector2
	GridCell GridCell
}

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

		// Coincident boids have no well-defined separation direction.
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

func BoidAction(b *Boid) {
	nearby := getNearbyBoids(b)

	cohesion := boidCohesion(b, nearby)
	separation := boidSeparation(b, nearby)
	alignment := boidAlignment(b, nearby)

	steering := rl.Vector2Add(
		rl.Vector2Scale(cohesion, CohesionValue),
		rl.Vector2Add(
			rl.Vector2Scale(separation, SeparationValue),
			rl.Vector2Scale(alignment, AlignmentValue),
		),
	)

	steering = limitVector(steering, MaxForce)

	dt := rl.GetFrameTime()

	b.Velocity = rl.Vector2Add(
		b.Velocity,
		rl.Vector2Scale(steering, dt),
	)
	b.Velocity = limitVector(b.Velocity, MaxSpeed)

	b.Position = rl.Vector2Add(
		b.Position,
		rl.Vector2Scale(b.Velocity, dt),
	)

	if b.Position.X < 0.0 {
		b.Position.X += Bounds.Width
	} else if b.Position.X > Bounds.Width {
		b.Position.X -= Bounds.Width
	}
	if b.Position.Y < 0.0 {
		b.Position.Y += Bounds.Height
	} else if b.Position.Y > Bounds.Height {
		b.Position.Y -= Bounds.Height
	}

	if DoGridOptimize {
		updateBoidGrid(b)
	}
}

func updateBoidGrid(b *Boid) {
	b.GridCell.X = floorDiv(b.Position.X, DetectionRange)
	b.GridCell.Y = floorDiv(b.Position.Y, DetectionRange)
}

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

var lineLength float32 = 24.0

func DrawBoid(b *Boid, highlight bool) {
	col := rl.Blue
	if highlight {
		col = rl.Red
	}
	vel := rl.Vector2Normalize(b.Velocity)
	if highlight {
		rl.DrawCircleLines(int32(b.Position.X), int32(b.Position.Y), DetectionRange, rl.LightGray)
	}
	rl.DrawLine(int32(b.Position.X), int32(b.Position.Y), int32(b.Position.X+vel.X*lineLength), int32(b.Position.Y+vel.Y*lineLength), rl.LightGray)
	rl.DrawPoly(
		b.Position,
		3,
		6,
		rl.Vector2Angle(right2d, b.Velocity)*rl.Rad2deg,
		col,
	)
}

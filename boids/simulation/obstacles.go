package simulation

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Obstacle represents a circular obstacle.
type Obstacle struct {
	Position rl.Vector2
	Radius   float32
}

// RectObstacle represents a rectangular obstacle.
type RectObstacle struct {
	Bounds rl.Rectangle
}

// CurrentMap controls which map is displayed.
// 1 = Obstacle Field
// 2 = Maze
var CurrentMap = 1

// ToggleMap switches between Map 1 and Map 2.
func ToggleMap() {
	if CurrentMap == 1 {
		CurrentMap = 2
	} else {
		CurrentMap = 1
	}
}

// GetMapName returns the name of the current map.
func GetMapName() string {
	if CurrentMap == 1 {
		return "MAP 1 - OBSTACLE FIELD"
	}

	return "MAP 2 - MAZE"
}

// ---------------- MAP 1 ----------------

// Circular obstacles.
var Obstacles = []Obstacle{
	{Position: rl.Vector2{X: 230, Y: 170}, Radius: 42},
	{Position: rl.Vector2{X: 430, Y: 190}, Radius: 32},
	{Position: rl.Vector2{X: 830, Y: 155}, Radius: 48},
	{Position: rl.Vector2{X: 930, Y: 390}, Radius: 42},
	{Position: rl.Vector2{X: 300, Y: 570}, Radius: 38},
	{Position: rl.Vector2{X: 710, Y: 590}, Radius: 34},
}

// Rectangular obstacles for Map 1.
var RectObstacles = []RectObstacle{
	{Bounds: rl.Rectangle{X: 550, Y: 100, Width: 150, Height: 35}},
	{Bounds: rl.Rectangle{X: 100, Y: 360, Width: 180, Height: 35}},
	{Bounds: rl.Rectangle{X: 650, Y: 420, Width: 140, Height: 40}},
}

// ---------------- MAP 2 ----------------

// Rectangular walls forming the maze.
// Gaps between the walls allow boids to pass through.
var MazeObstacles = []RectObstacle{
	// Upper section.
	{Bounds: rl.Rectangle{X: 150, Y: 120, Width: 330, Height: 25}},
	{Bounds: rl.Rectangle{X: 560, Y: 120, Width: 330, Height: 25}},

	{Bounds: rl.Rectangle{X: 150, Y: 120, Width: 25, Height: 190}},
	{Bounds: rl.Rectangle{X: 465, Y: 120, Width: 25, Height: 165}},
	{Bounds: rl.Rectangle{X: 875, Y: 120, Width: 25, Height: 190}},

	// Middle section.
	{Bounds: rl.Rectangle{X: 275, Y: 285, Width: 300, Height: 25}},
	{Bounds: rl.Rectangle{X: 635, Y: 285, Width: 240, Height: 25}},

	{Bounds: rl.Rectangle{X: 275, Y: 285, Width: 25, Height: 190}},
	{Bounds: rl.Rectangle{X: 600, Y: 285, Width: 25, Height: 165}},
	{Bounds: rl.Rectangle{X: 850, Y: 285, Width: 25, Height: 190}},

	// Lower section.
	{Bounds: rl.Rectangle{X: 275, Y: 450, Width: 325, Height: 25}},
	{Bounds: rl.Rectangle{X: 675, Y: 450, Width: 200, Height: 25}},

	{Bounds: rl.Rectangle{X: 450, Y: 450, Width: 25, Height: 125}},
	{Bounds: rl.Rectangle{X: 725, Y: 450, Width: 25, Height: 125}},

	// Bottom section.
	{Bounds: rl.Rectangle{X: 150, Y: 575, Width: 300, Height: 25}},
	{Bounds: rl.Rectangle{X: 500, Y: 575, Width: 350, Height: 25}},
}

// ---------------- DRAWING ----------------

// Draw obstacles belonging to the selected map.
func DrawObstacles() {
	if CurrentMap == 1 {
		drawMap1()
	} else {
		drawMap2()
	}
}

// Draw Map 1.
func drawMap1() {
	for _, obstacle := range Obstacles {
		rl.DrawCircleV(obstacle.Position, obstacle.Radius, rl.Gray)

		rl.DrawCircleLines(
			int32(obstacle.Position.X),
			int32(obstacle.Position.Y),
			obstacle.Radius,
			rl.DarkGray,
		)
	}

	for _, obstacle := range RectObstacles {
		r := obstacle.Bounds

		rl.DrawRectangleRec(r, rl.Gray)

		rl.DrawRectangleLines(
			int32(r.X),
			int32(r.Y),
			int32(r.Width),
			int32(r.Height),
			rl.DarkGray,
		)
	}
}

// Draw Map 2.
func drawMap2() {
	for _, obstacle := range MazeObstacles {
		r := obstacle.Bounds

		rl.DrawRectangleRec(r, rl.DarkGray)

		rl.DrawRectangleLines(
			int32(r.X),
			int32(r.Y),
			int32(r.Width),
			int32(r.Height),
			rl.Black,
		)
	}
}

// ---------------- AVOIDANCE ----------------

// Calculate avoidance for a circular obstacle.
func avoidCircle(b *Boid, position rl.Vector2, radius float32) rl.Vector2 {
	if b.Velocity.Length() == 0 {
		return rl.Vector2Zero()
	}

	forward := rl.Vector2Normalize(b.Velocity)
	lookAheadDistance := float32(180)
	extraClearance := float32(65)

	toObstacle := rl.Vector2Subtract(position, b.Position)
	projection := rl.Vector2DotProduct(toObstacle, forward)

	if projection < 0 || projection > lookAheadDistance {
		return rl.Vector2Zero()
	}

	closestPoint := rl.Vector2Add(
		b.Position,
		rl.Vector2Scale(forward, projection),
	)

	away := rl.Vector2Subtract(closestPoint, position)
	distance := away.Length()
	safeDistance := radius + extraClearance

	if distance >= safeDistance {
		return rl.Vector2Zero()
	}

	if distance < 0.001 {
		away = rl.Vector2{X: forward.Y, Y: -forward.X}
	} else {
		away = rl.Vector2Normalize(away)
	}

	urgency := (safeDistance - distance) / safeDistance

	desiredVelocity := rl.Vector2Scale(away, MaxSpeed)
	avoidForce := rl.Vector2Subtract(desiredVelocity, b.Velocity)

	return rl.Vector2Scale(avoidForce, urgency)
}

// Calculate avoidance for a rectangular wall.
// Unlike Map 1, this checks the actual rectangle instead of
// treating the entire wall as a large circle.
func avoidRectangle(b *Boid, bounds rl.Rectangle) rl.Vector2 {
	if b.Velocity.Length() == 0 {
		return rl.Vector2Zero()
	}

	forward := rl.Vector2Normalize(b.Velocity)

	// Direction perpendicular to the boid's movement.
	side := rl.Vector2{
		X: -forward.Y,
		Y: forward.X,
	}

	// Find the centre of the rectangle.
	center := rl.Vector2{
		X: bounds.X + bounds.Width/2,
		Y: bounds.Y + bounds.Height/2,
	}

	toWall := rl.Vector2Subtract(center, b.Position)

	// Distance along the boid's forward direction.
	projection := rl.Vector2DotProduct(toWall, forward)

	// Distance to the left or right of the boid.
	sideDistance := rl.Vector2DotProduct(toWall, side)

	// Calculate how far the rectangle extends in both directions.
	halfWidth := bounds.Width / 2
	halfHeight := bounds.Height / 2

	forwardExtent :=
		halfWidth*float32(math.Abs(float64(forward.X))) +
			halfHeight*float32(math.Abs(float64(forward.Y)))

	sideExtent :=
		halfWidth*float32(math.Abs(float64(side.X))) +
			halfHeight*float32(math.Abs(float64(side.Y)))

	lookAheadDistance := float32(180)
	extraClearance := float32(65)

	// The wall is either behind us or too far ahead.
	if projection+forwardExtent < 0 ||
		projection-forwardExtent > lookAheadDistance {
		return rl.Vector2Zero()
	}

	// The wall is not close enough to our movement path.
	distance := float32(math.Abs(float64(sideDistance)))
	safeDistance := sideExtent + extraClearance

	if distance >= safeDistance {
		return rl.Vector2Zero()
	}

	// Steer away from the wall.
	var away rl.Vector2

	if sideDistance > 0 {
		away = rl.Vector2Scale(side, -1)
	} else if sideDistance < 0 {
		away = side
	} else {
		away = rl.Vector2{
			X: forward.Y,
			Y: -forward.X,
		}
	}

	urgency := (safeDistance - distance) / safeDistance

	desiredVelocity := rl.Vector2Scale(away, MaxSpeed)
	avoidForce := rl.Vector2Subtract(desiredVelocity, b.Velocity)

	return rl.Vector2Scale(avoidForce, urgency)
}

// Calculate avoidance for all obstacles on the current map.
func obstacleAvoidance(b *Boid) rl.Vector2 {
	steering := rl.Vector2Zero()
	count := 0

	if CurrentMap == 1 {
		// Check circular obstacles.
		for _, obstacle := range Obstacles {
			force := avoidCircle(b, obstacle.Position, obstacle.Radius)

			if force.Length() > 0 {
				steering = rl.Vector2Add(steering, force)
				count++
			}
		}

		// Keep the original rectangular avoidance on Map 1.
		for _, obstacle := range RectObstacles {
			r := obstacle.Bounds

			center := rl.Vector2{
				X: r.X + r.Width/2,
				Y: r.Y + r.Height/2,
			}

			radius := float32(math.Sqrt(
				float64(r.Width*r.Width+r.Height*r.Height),
			)) / 2

			force := avoidCircle(b, center, radius)

			if force.Length() > 0 {
				steering = rl.Vector2Add(steering, force)
				count++
			}
		}
	} else {
		// Map 2 uses proper rectangular wall avoidance.
		for _, obstacle := range MazeObstacles {
			force := avoidRectangle(b, obstacle.Bounds)

			if force.Length() > 0 {
				steering = rl.Vector2Add(steering, force)
				count++
			}
		}
	}

	if count == 0 {
		return rl.Vector2Zero()
	}

	return rl.Vector2Scale(steering, 1.0/float32(count))
}

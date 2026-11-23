package simulation

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// The three simulation modes are free flocking, obstacles, and maze.
const (
	MapNone = iota
	MapObstacles
	MapMaze
)

// Obstacle defines a circular barrier by its centre and radius.
type Obstacle struct {
	Position rl.Vector2
	Radius   float32
}

// RectObstacle defines an axis-aligned wall by its bounds.
type RectObstacle struct {
	Bounds rl.Rectangle
}

// Start with no map, so the boids can flock freely.
var CurrentMap = MapNone

// ToggleMap cycles through the maps and moves boids clear of newly active obstacles.
func ToggleMap() {
	CurrentMap = (CurrentMap + 1) % 3

	// When a map appears, push boids out of any newly occupied walls.
	if CurrentMap != MapNone {
		for _, b := range Boids {
			ResolveObstacleCollision(b, b.Position)
		}
	}
}

// GetMapName returns the active map label shown in the interface.
func GetMapName() string {
	switch CurrentMap {
	case MapNone:
		return "NO MAP - FREE FLOCKING"
	case MapObstacles:
		return "MAP 1 - OBSTACLES"
	case MapMaze:
		return "MAP 2 - MAZE"
	default:
		return "UNKNOWN MAP"
	}
}

// Map 1 remains unchanged.
var Obstacles = []Obstacle{
	{Position: rl.Vector2{X: 230, Y: 170}, Radius: 42},
	{Position: rl.Vector2{X: 430, Y: 190}, Radius: 32},
	{Position: rl.Vector2{X: 830, Y: 155}, Radius: 48},
	{Position: rl.Vector2{X: 930, Y: 390}, Radius: 42},
	{Position: rl.Vector2{X: 300, Y: 570}, Radius: 38},
	{Position: rl.Vector2{X: 710, Y: 590}, Radius: 34},
}

var RectObstacles = []RectObstacle{
	{Bounds: rl.Rectangle{X: 550, Y: 100, Width: 150, Height: 35}},
	{Bounds: rl.Rectangle{X: 100, Y: 360, Width: 180, Height: 35}},
	{Bounds: rl.Rectangle{X: 650, Y: 420, Width: 140, Height: 40}},
}

// Map 2: a tidy, rectangular maze sized for the 1100x700 playfield.
//
// All walls are 24px thick. Connected walls MEET EDGE-TO-EDGE (rather
// than only touching at their corners), so there are no accidental gaps.
//
// Outer frame: top entrance 120px wide, bottom exit 120px wide.
// The two main bends have 190px+ openings and internal lanes are at
// least 106px wide, giving the boids room to navigate and turn.
var MazeObstacles = []RectObstacle{
	// Outer frame, with intentional entrances/exits.
	{Bounds: rl.Rectangle{X: 95, Y: 65, Width: 350, Height: 24}},   // top left
	{Bounds: rl.Rectangle{X: 565, Y: 65, Width: 440, Height: 24}},  // top right
	{Bounds: rl.Rectangle{X: 95, Y: 89, Width: 24, Height: 522}},   // left side
	{Bounds: rl.Rectangle{X: 981, Y: 89, Width: 24, Height: 522}},  // right side
	{Bounds: rl.Rectangle{X: 95, Y: 611, Width: 510, Height: 24}},  // bottom left
	{Bounds: rl.Rectangle{X: 725, Y: 611, Width: 280, Height: 24}}, // bottom right

	// Two long, alternating barriers create a wide winding route.
	{Bounds: rl.Rectangle{X: 119, Y: 240, Width: 670, Height: 24}}, // right opening: 192px
	{Bounds: rl.Rectangle{X: 310, Y: 434, Width: 671, Height: 24}}, // left opening: 191px

	// Short right-angle branches add interest without blocking passages.
	// All branches join flush with their parent wall.
	{Bounds: rl.Rectangle{X: 260, Y: 89, Width: 24, Height: 40}},
	{Bounds: rl.Rectangle{X: 690, Y: 89, Width: 24, Height: 40}},
	{Bounds: rl.Rectangle{X: 380, Y: 264, Width: 24, Height: 55}},
	{Bounds: rl.Rectangle{X: 600, Y: 264, Width: 24, Height: 50}},
	{Bounds: rl.Rectangle{X: 500, Y: 382, Width: 24, Height: 52}},
	{Bounds: rl.Rectangle{X: 770, Y: 382, Width: 24, Height: 52}},
	{Bounds: rl.Rectangle{X: 365, Y: 564, Width: 24, Height: 47}},
	{Bounds: rl.Rectangle{X: 850, Y: 564, Width: 24, Height: 47}},
}

// activeRects returns the rectangular barriers for the active map.
func activeRects() []RectObstacle {
	switch CurrentMap {
	case MapObstacles:
		return RectObstacles
	case MapMaze:
		return MazeObstacles
	default:
		return nil
	}
}

// DrawObstacles draws the barriers for the active map.
func DrawObstacles() {
	switch CurrentMap {
	case MapNone:
		return
	case MapObstacles:
		for _, o := range Obstacles {
			rl.DrawCircleV(o.Position, o.Radius, rl.Gray)
			rl.DrawCircleLines(int32(o.Position.X), int32(o.Position.Y), o.Radius, rl.DarkGray)
		}
		for _, o := range RectObstacles {
			rl.DrawRectangleRec(o.Bounds, rl.Gray)
			r := o.Bounds
			rl.DrawRectangleLines(int32(r.X), int32(r.Y), int32(r.Width), int32(r.Height), rl.Black)
		}
	case MapMaze:
		drawMaze()
	}
}

// drawMaze draws the maze walls and marks the entrance and exit.
func drawMaze() {
	// Draw ALL shadows first, then ALL walls, so connected joints look solid.
	shadow := rl.Color{R: 15, G: 23, B: 38, A: 165}
	wall := rl.Color{R: 65, G: 86, B: 115, A: 255}
	edge := rl.Color{R: 99, G: 125, B: 157, A: 255}

	for _, o := range MazeObstacles {
		r := o.Bounds
		rl.DrawRectangleRec(rl.Rectangle{
			X: r.X + 4, Y: r.Y + 5, Width: r.Width, Height: r.Height,
		}, shadow)
	}
	for _, o := range MazeObstacles {
		rl.DrawRectangleRec(o.Bounds, wall)
	}
	// A restrained 2px highlight gives the walls some depth.
	for _, o := range MazeObstacles {
		r := o.Bounds
		rl.DrawRectangleRec(rl.Rectangle{
			X: r.X, Y: r.Y, Width: r.Width, Height: 2,
		}, edge)
	}

	// Subtle labels indicate the two deliberate gaps in the frame.
	rl.DrawText("IN", 493, 45, 15, rl.LightGray)
	rl.DrawText("OUT", 647, 645, 15, rl.LightGray)
}

// clamp keeps a value between the given limits.
func clamp(x, lo, hi float32) float32 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// dot returns the dot product of two vectors.
func dot(a, b rl.Vector2) float32 { return a.X*b.X + a.Y*b.Y }

// sub subtracts the second vector from the first.
func sub(a, b rl.Vector2) rl.Vector2 { return rl.Vector2Subtract(a, b) }

// add returns the sum of two vectors.
func add(a, b rl.Vector2) rl.Vector2 { return rl.Vector2Add(a, b) }

// mul scales a vector by the given factor.
func mul(a rl.Vector2, k float32) rl.Vector2 { return rl.Vector2Scale(a, k) }

// length returns the magnitude of a vector.
func length(a rl.Vector2) float32 { return a.Length() }

// unit returns a direction of length one, using the positive X axis for tiny vectors.
func unit(a rl.Vector2) rl.Vector2 {
	if length(a) < 0.0001 {
		return rl.Vector2{X: 1}
	}
	return rl.Vector2Normalize(a)
}

// closestRect finds the closest point on or inside the rectangle.
func closestRect(p rl.Vector2, r rl.Rectangle) rl.Vector2 {
	return rl.Vector2{X: clamp(p.X, r.X, r.X+r.Width), Y: clamp(p.Y, r.Y, r.Y+r.Height)}
}

// circleSurface returns the outward normal and signed distance to a circle.
func circleSurface(p rl.Vector2, o Obstacle) (rl.Vector2, float32) {
	v := sub(p, o.Position)
	d := length(v)
	return unit(v), d - o.Radius
}

// rectSurface returns the outward normal and signed distance to a rectangle.
func rectSurface(p rl.Vector2, r rl.Rectangle) (rl.Vector2, float32) {
	q := closestRect(p, r)
	v := sub(p, q)
	d := length(v)
	if d > 0.0001 {
		return mul(v, 1/d), d
	}

	left := p.X - r.X
	right := r.X + r.Width - p.X
	top := p.Y - r.Y
	bottom := r.Y + r.Height - p.Y
	m := left
	n := rl.Vector2{X: -1}
	if right < m {
		m = right
		n = rl.Vector2{X: 1}
	}
	if top < m {
		m = top
		n = rl.Vector2{Y: -1}
	}
	if bottom < m {
		m = bottom
		n = rl.Vector2{Y: 1}
	}
	return n, -m
}

// nearestSurface finds the nearest active obstacle surface and its outward normal.
func nearestSurface(p rl.Vector2) (rl.Vector2, float32) {
	if CurrentMap == MapNone {
		return rl.Vector2Zero(), float32(math.MaxFloat32)
	}

	best := float32(math.MaxFloat32)
	normal := rl.Vector2{X: 1}
	if CurrentMap == MapObstacles {
		for _, o := range Obstacles {
			n, d := circleSurface(p, o)
			if d < best {
				best, normal = d, n
			}
		}
	}
	for _, o := range activeRects() {
		n, d := rectSurface(p, o.Bounds)
		if d < best {
			best, normal = d, n
		}
	}
	return normal, best
}

// mazeCornerDirection finds an exit direction when two nearby wall faces meet.
func mazeCornerDirection(p rl.Vector2) rl.Vector2 {
	const reach = 22
	var outward rl.Vector2
	for _, wall := range MazeObstacles {
		r := wall.Bounds
		// Only count faces alongside the boid. A lone wall tip has space on both sides.
		if p.Y >= r.Y && p.Y <= r.Y+r.Height {
			if gap := p.X - (r.X + r.Width); gap >= 0 && gap < reach {
				outward.X = 1
			} else if gap := r.X - p.X; gap >= 0 && gap < reach {
				outward.X = -1
			}
		}
		if p.X >= r.X && p.X <= r.X+r.Width {
			if gap := p.Y - (r.Y + r.Height); gap >= 0 && gap < reach {
				outward.Y = 1
			} else if gap := r.Y - p.Y; gap >= 0 && gap < reach {
				outward.Y = -1
			}
		}
	}
	if outward.X != 0 && outward.Y != 0 {
		return unit(outward)
	}
	return rl.Vector2Zero()
}

// obstacleAvoidance looks ahead and steers the boid around nearby obstacles.
func obstacleAvoidance(b *Boid) rl.Vector2 {
	if CurrentMap == MapNone {
		return rl.Vector2Zero()
	}

	// Two walls can cancel both velocity components, so check corners even at rest.
	if CurrentMap == MapMaze {
		escape := mazeCornerDirection(b.Position)
		if escape != (rl.Vector2{}) {
			return sub(mul(escape, MaxSpeed), b.Velocity)
		}
	}

	speed := length(b.Velocity)
	if speed < 0.001 {
		return rl.Vector2Zero()
	}
	dir := unit(b.Velocity)
	look := clamp(speed*0.18, 35, 115)
	best := float32(math.MaxFloat32)
	normal := rl.Vector2Zero()
	at := float32(0)
	for i := 0; i <= 12; i++ {
		t := float32(i) / 12
		point := add(b.Position, mul(dir, look*t))
		n, d := nearestSurface(point)
		if d < best {
			best, normal, at = d, n, t
		}
	}
	if best > 22 {
		return rl.Vector2Zero()
	}

	// Slide alongside a wall instead of always turning straight backwards.
	tangent := rl.Vector2{X: -normal.Y, Y: normal.X}
	if dot(tangent, dir) < 0 {
		tangent = mul(tangent, -1)
	}
	urgency := clamp((22-best)/42, 0, 1)
	desired := add(mul(tangent, 0.75), mul(normal, 0.25+0.6*urgency))
	if at == 0 && best < 8 {
		desired = add(desired, mul(normal, 1.2))
	}
	return mul(sub(mul(unit(desired), MaxSpeed), b.Velocity), urgency)
}

// ResolveObstacleCollision moves a boid along its swept path and removes motion into walls.
func ResolveObstacleCollision(b *Boid, previous rl.Vector2) {
	if CurrentMap == MapNone {
		return
	}
	const margin float32 = 7
	delta := sub(b.Position, previous)
	distance := length(delta)
	steps := int(math.Ceil(float64(distance / 4)))
	if steps < 1 {
		steps = 1
	}
	if steps > 256 {
		steps = 256
	}
	pos := previous

	// A map switch can create a wall at a boid's old position.
	for j := 0; j < 8; j++ {
		n, d := nearestSurface(pos)
		if d >= margin {
			break
		}
		pos = add(pos, mul(n, margin-d+0.1))
		vn := dot(b.Velocity, n)
		if vn < 0 {
			b.Velocity = sub(b.Velocity, mul(n, vn))
		}
	}

	step := mul(delta, 1/float32(steps))
	for i := 0; i < steps; i++ {
		candidate := add(pos, step)
		for j := 0; j < 4; j++ {
			n, d := nearestSurface(candidate)
			if d >= margin {
				break
			}
			candidate = add(candidate, mul(n, margin-d+0.1))
			vn := dot(b.Velocity, n)
			if vn < 0 {
				b.Velocity = sub(b.Velocity, mul(n, vn))
			}
			into := dot(step, n)
			if into < 0 {
				step = sub(step, mul(n, into))
			}
		}
		pos = candidate
	}
	b.Position = pos
}

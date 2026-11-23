package tests

import (
	"boids/simulation"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// setupObstacleSettings saves and restores the active map and flock settings.
func setupObstacleSettings(t *testing.T) {
	setupBoidSettings(t)
	oldMap := simulation.CurrentMap
	oldForce := simulation.MaxForce
	oldWorkers := simulation.Workers
	t.Cleanup(func() {
		simulation.CurrentMap = oldMap
		simulation.MaxForce = oldForce
		simulation.Workers = oldWorkers
	})
	simulation.CurrentMap = simulation.MapNone
	simulation.MaxForce = 200
}

// TestMapSwitch checks that all three map modes are selected in order.
func TestMapSwitch(t *testing.T) {
	setupObstacleSettings(t)
	for _, expected := range []int{simulation.MapObstacles, simulation.MapMaze, simulation.MapNone} {
		simulation.ToggleMap()
		if simulation.CurrentMap != expected {
			t.Fatalf("map is %d; expected %d", simulation.CurrentMap, expected)
		}
	}
}

// TestNoMapMovement checks that free flocking uses the original movement.
func TestNoMapMovement(t *testing.T) {
	setupObstacleSettings(t)
	boid := &simulation.Boid{Position: rl.Vector2{X: 100, Y: 100}, Velocity: rl.Vector2{X: 40}}
	simulation.Boids = []*simulation.Boid{boid}
	simulation.BoidAction(boid, 0.5)
	if boid.Position != (rl.Vector2{X: 120, Y: 100}) || boid.Velocity != (rl.Vector2{X: 40}) {
		t.Fatalf("unexpected free flocking movement: %+v", boid)
	}
}

// TestObstacleCollision checks that a boid cannot jump through a wall or circle.
func TestObstacleCollision(t *testing.T) {
	setupObstacleSettings(t)
	// Each case gives a swept path and the furthest X allowed by the obstacle.
	cases := []struct {
		name       string
		mode       int
		start, end rl.Vector2
		maxX       float32
	}{
		{name: "circle", mode: simulation.MapObstacles, start: rl.Vector2{X: 150, Y: 170}, end: rl.Vector2{X: 310, Y: 170}, maxX: 181},
		{name: "rectangle", mode: simulation.MapObstacles, start: rl.Vector2{X: 50, Y: 375}, end: rl.Vector2{X: 330, Y: 375}, maxX: 93},
		{name: "maze wall", mode: simulation.MapMaze, start: rl.Vector2{X: 50, Y: 200}, end: rl.Vector2{X: 150, Y: 200}, maxX: 88},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			simulation.CurrentMap = test.mode
			boid := &simulation.Boid{Position: test.end, Velocity: rl.Vector2{X: 100}}
			simulation.ResolveObstacleCollision(boid, test.start)
			if boid.Position.X > test.maxX {
				t.Fatalf("boid crossed obstacle: %+v", boid.Position)
			}
			if boid.Velocity.X > 0 {
				t.Fatalf("boid still moving into obstacle: %+v", boid.Velocity)
			}
		})
	}
}

// TestMazeCornerEscape checks that a boid can leave a two-wall junction after impact.
func TestMazeCornerEscape(t *testing.T) {
	setupObstacleSettings(t)
	simulation.CurrentMap = simulation.MapMaze
	// Each case starts seven pixels from both walls, with a direction into open space.
	cases := []struct {
		name           string
		start, outward rl.Vector2
	}{
		{"upper left", rl.Vector2{X: 126, Y: 96}, rl.Vector2{X: 1, Y: 1}},
		{"upper right", rl.Vector2{X: 974, Y: 96}, rl.Vector2{X: -1, Y: 1}},
		{"below left barrier", rl.Vector2{X: 126, Y: 271}, rl.Vector2{X: 1, Y: 1}},
		{"above right barrier", rl.Vector2{X: 974, Y: 427}, rl.Vector2{X: -1, Y: -1}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			previous := rl.Vector2Add(test.start, rl.Vector2Scale(test.outward, 4))
			boid := &simulation.Boid{
				Position: rl.Vector2Subtract(test.start, rl.Vector2Scale(test.outward, 4)),
				Velocity: rl.Vector2Scale(test.outward, -70),
			}
			simulation.ResolveObstacleCollision(boid, previous)
			if boid.Velocity != (rl.Vector2{}) {
				t.Fatalf("expected impact with both walls, got velocity %+v", boid.Velocity)
			}
			simulation.Boids = []*simulation.Boid{boid}
			for range 90 {
				simulation.BoidAction(boid, 1.0/60)
				offset := rl.Vector2Subtract(boid.Position, test.start)
				if offset.X*test.outward.X < -0.01 || offset.Y*test.outward.Y < -0.01 {
					t.Fatalf("boid entered a corner wall at %+v", boid.Position)
				}
				if offset.X*test.outward.X > 18 && offset.Y*test.outward.Y > 18 {
					return
				}
			}
			t.Fatalf("boid stayed near the corner at %+v", boid.Position)
		})
	}
}

// TestMapSwitchGrid checks that boids pushed out of a wall update their grid cell.
func TestMapSwitchGrid(t *testing.T) {
	setupObstacleSettings(t)
	simulation.DoGridOptimize = true
	boid := &simulation.Boid{Position: rl.Vector2{X: 230, Y: 170}}
	simulation.Boids = []*simulation.Boid{boid}
	simulation.RefreshBoidGrid()
	simulation.ToggleMap()
	simulation.RefreshBoidGrid()
	if rl.Vector2Distance(boid.Position, simulation.Obstacles[0].Position) < simulation.Obstacles[0].Radius+7 {
		t.Fatal("boid was not pushed out of the circle")
	}
	if boid.GridCell.X != int(boid.Position.X/simulation.DetectionRange) || boid.GridCell.Y != int(boid.Position.Y/simulation.DetectionRange) {
		t.Fatal("boid grid cell was not refreshed")
	}
}

// TestSerialParallelObstacleMovement checks that both update methods use the same obstacle steering.
func TestSerialParallelObstacleMovement(t *testing.T) {
	setupObstacleSettings(t)
	simulation.CurrentMap = simulation.MapObstacles
	simulation.Workers = 2
	startingBoid := simulation.Boid{Position: rl.Vector2{X: 150, Y: 170}, Velocity: rl.Vector2{X: 100}}
	serial := startingBoid
	simulation.Boids = []*simulation.Boid{&serial}
	simulation.BoidAction(&serial, 0.1)
	parallel := startingBoid
	simulation.Boids = []*simulation.Boid{&parallel}
	simulation.UpdateBoids(0.1)
	if serial != parallel {
		t.Fatalf("serial state %+v differs from parallel state %+v", serial, parallel)
	}
	if serial.Velocity == startingBoid.Velocity {
		t.Fatal("obstacle avoidance did not change the velocity")
	}
}

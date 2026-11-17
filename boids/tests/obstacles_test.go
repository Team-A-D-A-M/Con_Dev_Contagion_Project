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

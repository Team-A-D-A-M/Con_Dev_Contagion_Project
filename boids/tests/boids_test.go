package tests

import (
	"boids/simulation"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// setupBoidSettings sets known values and restores them after the test.
func setupBoidSettings(t *testing.T) {
	// Save the current values of the variables declared in simulation/parameters.go.
	oldBoids := simulation.Boids
	oldGrid := simulation.DoGridOptimize
	oldRange := simulation.DetectionRange
	oldSeparation := simulation.DesiredSeparation
	oldSpeed := simulation.MaxSpeed

	// After the test finishes, restore the values saved above.
	t.Cleanup(func() {
		simulation.Boids = oldBoids
		simulation.DoGridOptimize = oldGrid
		simulation.DetectionRange = oldRange
		simulation.DesiredSeparation = oldSeparation
		simulation.MaxSpeed = oldSpeed
	})

	// Temporarily replace the current flock and settings with values for this test.
	simulation.Boids = []*simulation.Boid{}
	simulation.DoGridOptimize = false
	simulation.DetectionRange = 80
	simulation.DesiredSeparation = 24
	simulation.MaxSpeed = 100
}

// TestDetectionRange checks boids inside, at and outside the range.
func TestDetectionRange(t *testing.T) {
	// Create a list of test cases
	cases := []struct {
		name               string  // Name shown in the test results
		x                  float32 // The second boid's X position.
		expectedNeighbours int     // How many neighbours should be detected.
	}{
		{name: "inside", x: 150, expectedNeighbours: 1},           // 50 units away from first boid
		{name: "exactly at range", x: 180, expectedNeighbours: 0}, // 80 units away
		{name: "outside", x: 300, expectedNeighbours: 0},          // 200 units away
	}

	// run for each case
	for _, test := range cases {
		// Use the case's name so we can see which one passes or fails.
		t.Run(test.name, func(t *testing.T) {
			setupBoidSettings(t)

			// The first boid is at X = 100. Detection range is 80.
			boid := &simulation.Boid{
				Position: rl.Vector2{X: 100, Y: 100},
			}
			other := &simulation.Boid{
				Position: rl.Vector2{X: test.x, Y: 100},
			}

			// Put both boids in the flock.
			simulation.Boids = []*simulation.Boid{boid, other}

			// Find the first boid's neighbours.
			nearby := simulation.GetNearbyBoids(boid)

			// Fail and stop this case if the neighbour count is wrong.
			if len(nearby) != test.expectedNeighbours {
				t.Fatalf("found %d neighbours; expectedNeighbours %d",
					len(nearby), test.expectedNeighbours)
			}

			// If we expected a neighbour, check that it is the second boid.
			if test.expectedNeighbours == 1 && nearby[0].Boid != other {
				t.Error("detected the wrong boid")
			}
		})
	}
}

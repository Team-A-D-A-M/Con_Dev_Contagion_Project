package tests

import (
	"boids/simulation"
	"math/rand/v2"
	"strconv"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// BenchmarkSerial measures serial updates for different numbers of boids.
// Each update starts from the same saved state.
// Creating, resetting and drawing the boids are not included in the timing.
func BenchmarkSerial(b *testing.B) {

	amounts := []int{100, 500, 1000}

	for _, amount := range amounts {
		// Convert the boid count to text and run a benchmark with that name.
		b.Run(strconv.Itoa(amount), func(b *testing.B) {
			// Create empty list of boid pointers
			simulation.Boids = []*simulation.Boid{}
			// Turn off grid optimisation for this first benchmark.
			simulation.DoGridOptimize = false

			// Use the same random sequence each time the benchmark runs.
			random := rand.New(rand.NewPCG(1, 2))

			// Create chosen number of boids
			for range amount {
				simulation.Boids = append(simulation.Boids, &simulation.Boid{
					// Start in a random spot within 150 pixels right and down from the centre.
					Position: rl.Vector2{
						X: random.Float32()*150 + simulation.Bounds.Width*0.5,
						Y: random.Float32()*150 + simulation.Bounds.Height*0.5},
				})
			}

			// Create a list to hold each boid's starting data.
			startingBoids := make([]simulation.Boid, len(simulation.Boids))

			// Save each boid's starting data.
			for i, boid := range simulation.Boids {
				startingBoids[i] = *boid
			}

			// Measure updates of all boids.
			for b.Loop() {
				// Pause timing while resetting the birds.
				b.StopTimer()

				// Copy the starting data of each boid back into each boid.
				for i, boid := range simulation.Boids {
					*boid = startingBoids[i]
				}

				// Time only the serial update.
				b.StartTimer()

				// Update each boid one at a time (serially).
				for _, boid := range simulation.Boids {
					simulation.BoidAction(boid, 1.0/60.0)
				}
			}
		})
	}
}

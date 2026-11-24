package tests

import (
	"boids/simulation"
	"math/rand/v2"
	"strconv"
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var boidCounts = []int{25, 50, 100, 500, 1000}

// BenchmarkSerial measures serial updates for different numbers of boids.
// Each update starts from the same saved state.
// Creating, resetting and drawing the boids are not included in the timing.
func BenchmarkSerial(b *testing.B) {
	for _, amount := range boidCounts {
		// Convert the boid count to text and run a benchmark with that name.
		b.Run(strconv.Itoa(amount), func(b *testing.B) {
			benchmarkFlock(b, amount, func() {
				// Update each boid one at a time (serially).
				for _, boid := range simulation.Boids {
					simulation.BoidAction(boid, 1.0/60.0)
				}
			})
		})
	}
}

// BenchmarkParallel measures concurrent updates for different numbers of boids and worker counts
// Each update starts from the same saved state
// Creating, resetting and drawing the boids are not included in the timing
// TODO Run the parallel benchmark, record the results, run the Serial benchmark (update with the same amount values) and compare the results
func BenchmarkParallel(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8, 10}

	for _, amount := range boidCounts {
		// Convert the boid count to text and run a benchmark with that name.
		b.Run(strconv.Itoa(amount), func(b *testing.B) {
			for _, workers := range workerCounts {
				// Run a separate benchmark for each worker count.
				b.Run(strconv.Itoa(workers)+"_workers", func(b *testing.B) {
					// Set the number of workers for this benchmark.
					simulation.Workers = workers

					benchmarkFlock(b, amount, func() {
						// Update the whole flock using the chosen workers.
						simulation.UpdateBoids(1.0 / 60.0)
					})
				})
			}
		})
	}
}

// benchmarkFlock sets up the flock and times the update passed in by either benchmark.
func benchmarkFlock(b *testing.B, amount int, update func()) {
	// Use a map without obstacles and turn off grid optimisation.
	simulation.CurrentMap = simulation.MapNone
	simulation.DoGridOptimize = false

	// Use the same starting positions in both benchmarks.
	random := rand.New(rand.NewPCG(1, 2))
	simulation.Boids = make([]*simulation.Boid, amount)
	startingBoids := make([]simulation.Boid, amount)

	// Create one boid for each slot in the flock.
	for i := range startingBoids {
		// Pick a starting position up to 150 pixels right and down from the centre.
		boid := &simulation.Boid{
			Position: rl.Vector2{
				X: random.Float32()*150 + simulation.Bounds.Width*0.5,
				Y: random.Float32()*150 + simulation.Bounds.Height*0.5,
			},
		}

		// Store the pointer to the boid that the benchmark will update.
		simulation.Boids[i] = boid

		// Copy the boid's starting data so we can reset it before each update.
		startingBoids[i] = *boid
	}

	// Repeat the benchmark. Each round measures one update of the whole flock.
	for b.Loop() {
		// Pause the timer so resetting the boids doesn't count.
		b.StopTimer()

		// Restore each boid's saved position and velocity so every round starts the same way.
		for i, boid := range simulation.Boids {
			*boid = startingBoids[i]
		}

		// Resume the timer and run one serial or parallel update.
		b.StartTimer()
		update()
	}
}

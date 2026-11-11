package simulation

import (
	"sync"
)

// Concurrent UpdateBoids for updating the population
func UpdateBoids(dt float32) {
	// make a snapshot array of boids with the same length as Boids and copy each *Boid value
	var snapshot []Boid
	snapshot = make([]Boid, len(Boids))

	for i, boid := range Boids {
		snapshot[i] = *boid
	}

	// each slot written once, based on the matching snapshot entry
	var nextStates = make([]Boid, len(snapshot))

	// make sure there aren't more workers than boids in the snapshot
	var workerCount = Workers
	if workerCount > len(snapshot) {
		workerCount = len(snapshot)
	}

	// queue of boid indexes workers take from to know which boid to update
	var jobs = make(chan int, len(snapshot))
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		// launch go routine for each worker
		go func() {
			defer wg.Done()
			// workers wait for jobs (blocks until index available)
			for i := range jobs {
				// calculate boids next state, just store result don't update yet
				nextStates[i] = BoidActionConcurrent(i, snapshot, dt)
			}
		}()
	}

	// puts every boid index into the jobs queue (unblocks the workers)
	for i := range snapshot {
		jobs <- i
	}
	// tells workers no more indexes coming
	close(jobs)
	// wait for workers to finish until main goroutine continues
	wg.Wait()

	// once all go routines have completed, update Boids with the next state
	for i := range Boids {
		*Boids[i] = nextStates[i]
	}
}

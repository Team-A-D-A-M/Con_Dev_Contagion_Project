package simulation

import "sync"

// Concurrent UpdateBoids for updating the population
func UpdateBoids(dt float32) {
	// make a snapshot array of boids with the same length as Boids and copy each *Boid value
	var snapshot []Boid
	snapshot = make([]Boid, len(Boids))

	for i, boid := range Boids {
		snapshot[i] = *boid
	}

	// each slot written once, based on the matching snapshot entry
	var nextStates []Boid
	nextStates = make([]Boid, len(snapshot))
	var wg sync.WaitGroup

	for i := range snapshot {
		wg.Add(1)
		// 1 go routine per boid
		go func(i int) {
			defer wg.Done()
			nextStates[i] = BoidActionConcurrent(i, snapshot, dt)
		}(i)
	}

	wg.Wait()

	// once all go routines have completed, update Boids with the next state
	for i := range Boids {
		*Boids[i] = nextStates[i]
	}
}

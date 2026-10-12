// author: Adam Noonan

package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// agent struct
type Agent struct {
	position        rl.Vector2
	infected        bool
	infectedCounter int
	resistance      float32
}

// move an agent a random amount
func agentMove(agent *Agent) {

	// generate random direction vector
	randomDirection := rl.Vector2{
		X: rand.Float32()*2 - 1,
		Y: rand.Float32()*2 - 1,
	}

	// set clamped X position
	agent.position.X = rl.Clamp(
		agent.position.X+randomDirection.X,
		bounds.X,
		bounds.X+bounds.Width,
	)

	// set clamped Y position
	agent.position.Y = rl.Clamp(
		agent.position.Y+randomDirection.Y,
		bounds.Y,
		bounds.X+bounds.Height,
	)

	// decrement agents infection counter
	agentDecrementInfect(agent)
}

func agentInfect(agent *Agent) {
	infectRoll := rand.Float32()

	if !agent.infected && infectRoll > agent.resistance {
		agent.infected = true
		agent.infectedCounter = 50
		agent.resistance *= 1.5
	}
}

func agentDecrementInfect(agent *Agent) {
	if agent.infected {
		agent.infectedCounter--
		if agent.infectedCounter == 0 {
			agent.infected = false
		}
	}
}

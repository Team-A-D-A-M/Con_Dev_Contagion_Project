package condevcontagionproject

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Agent struct {
	position        rl.Vector2
	infected        bool
	infectedCounter int
	resistance      float32
}

func agentMove(agent *Agent) {
	randomDirection := rl.Vector2{
		X: rand.Float32()*2 - 1,
		Y: rand.Float32()*2 - 1,
	}
	agent.position.X = rl.Clamp(
		agent.position.X+randomDirection.X,
		bounds.X, bounds.X+bounds.Width)
	agent.position.Y = rl.Clamp(
		agent.position.Y+randomDirection.Y,
		bounds.Y, bounds.X+bounds.Height)

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

package handlers

import (
	"sync"
	"time"
)

var (
	globalHub *clientHub
	gameState *GameState
	gameMu    sync.Mutex
)

func init() {
	globalHub = &clientHub{clients: make(map[uint64]chan []byte)}
	gameState = newGame()
	go runGameLoop()
}

func newGame() *GameState {
	return &GameState{
		Speed:     baseSpeed,
		Obstacles: []Obstacle{{X: float64(svgWidth) + 120, Width: 22, Height: 45}},
	}
}

func runGameLoop() {
	ticker := time.NewTicker(33 * time.Millisecond) // ~30 fps
	defer ticker.Stop()
	for range ticker.C {
		tick()
	}
}

func tick() {
	gameMu.Lock()
	gs := gameState
	if gs.Dead {
		gs.DeadTimer--
		if gs.DeadTimer <= 0 {
			gameState = newGame()
			gs = gameState
		}
	} else {
		applyGravity(gs)
		moveObstacles(gs)
		maybeSpawn(gs)
		gs.Score++
		gs.Speed = baseSpeed + float64(gs.Score)/400.0
		checkCollision(gs)
	}
	frame := RenderFrame(gs, globalHub.PlayerCount())
	gameMu.Unlock()
	globalHub.Broadcast(frame)
}

package handlers

import (
	"math/rand"
	"sync"
	"time"
)

// SVG canvas dimensions and game constants.
const (
	svgWidth    = 800
	svgHeight   = 200
	groundY     = 160 // SVG y-coordinate of the ground line

	dinoX      = 80
	dinoWidth  = 28
	dinoHeight = 40

	jumpImpulse = 11.0
	gravity     = 0.55
	baseSpeed   = 4.0
)

// Obstacle is a cactus sitting on the ground.
type Obstacle struct {
	X      float64
	Width  float64
	Height float64
}

// GameState holds all mutable game data. Owned by gameMu.
type GameState struct {
	DinoY     float64 // distance above ground (0 = standing)
	DinoVY    float64 // vertical velocity (+ve = upward)
	Obstacles []Obstacle
	Score     int
	Speed     float64
	Dead      bool
	DeadTimer int // ticks until auto-reset after death
}

// clientHub manages all connected SSE clients.
type clientHub struct {
	mu      sync.Mutex
	clients map[uint64]chan []byte
	nextID  uint64
	count   int
}

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
		Speed: baseSpeed,
		Obstacles: []Obstacle{
			{X: float64(svgWidth) + 120, Width: 22, Height: 45},
		},
	}
}

// Register adds a new SSE client. Returns its unique ID and the frame channel.
func (h *clientHub) Register() (uint64, chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	id := h.nextID
	h.nextID++
	ch := make(chan []byte, 2) // buffer of 2 to absorb brief slowness
	h.clients[id] = ch
	h.count++
	return id, ch
}

// Unregister removes a client and closes its channel.
func (h *clientHub) Unregister(id uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.clients[id]; ok {
		close(ch)
		delete(h.clients, id)
		h.count--
	}
}

// PlayerCount returns the number of currently connected players.
// Safe to call while gameMu is held (lock order: gameMu → hubMu).
func (h *clientHub) PlayerCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.count
}

// Broadcast sends a rendered frame to every connected client.
// Slow clients are silently dropped (non-blocking send).
func (h *clientHub) Broadcast(frame []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range h.clients {
		select {
		case ch <- frame:
		default: // drop frame for a lagging client
		}
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
		// --- Physics ---
		gs.DinoVY -= gravity
		gs.DinoY += gs.DinoVY
		if gs.DinoY < 0 {
			gs.DinoY = 0
			gs.DinoVY = 0
		}

		// Move obstacles leftward.
		for i := range gs.Obstacles {
			gs.Obstacles[i].X -= gs.Speed
		}

		// Remove obstacles that have scrolled off the left edge.
		var alive []Obstacle
		for _, o := range gs.Obstacles {
			if o.X+o.Width > 0 {
				alive = append(alive, o)
			}
		}
		gs.Obstacles = alive

		// Spawn a new obstacle once the rightmost one has passed x=450.
		needSpawn := true
		for _, o := range gs.Obstacles {
			if o.X > float64(svgWidth)-350 {
				needSpawn = false
				break
			}
		}
		if needSpawn {
			gs.Obstacles = append(gs.Obstacles, Obstacle{
				X:      float64(svgWidth) + float64(rand.Intn(220)),
				Width:  float64(18 + rand.Intn(18)),
				Height: float64(30 + rand.Intn(45)),
			})
		}

		// Score and speed.
		gs.Score++
		gs.Speed = baseSpeed + float64(gs.Score)/400.0

		// --- AABB Collision ---
		dinoL := float64(dinoX) + 2   // 2px inset for fairness
		dinoR := float64(dinoX+dinoWidth) - 2
		dinoB := float64(groundY) - gs.DinoY
		dinoT := dinoB - dinoHeight

		for _, o := range gs.Obstacles {
			obsL := o.X
			obsR := o.X + o.Width
			obsB := float64(groundY)
			obsT := obsB - o.Height

			if dinoR > obsL && dinoL < obsR && dinoB > obsT && dinoT < obsB {
				gs.Dead = true
				gs.DeadTimer = 125 // ~2 s at 60 fps
				break
			}
		}
	}

	// Render while gameMu is held so state is consistent.
	// hubMu is acquired inside PlayerCount() — safe: lock order gameMu→hubMu never reversed.
	playerCount := globalHub.PlayerCount()
	frame := RenderFrame(gs, playerCount)
	gameMu.Unlock()

	globalHub.Broadcast(frame)
}

// Jump sets an upward velocity on the dino if it is currently on the ground.
// Any connected player pressing the button affects the shared dino.
func Jump() {
	gameMu.Lock()
	defer gameMu.Unlock()
	if !gameState.Dead && gameState.DinoY == 0 {
		gameState.DinoVY = jumpImpulse
	}
}

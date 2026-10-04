package game

import (
	"math/rand"
	"sync"
	"time"
)

const (
	// Layout constants (pixels, matching CSS)
	groundY     = 200.0 // y at which dino stands (top of ground)
	dinoHeight  = 55.0
	jumpVel     = -14.0 // initial upward velocity (px/tick)
	gravity     = 1.0   // px/tick² downward
	arenaWidth  = 800.0
	maxSlots    = 4

	// Tick
	tickInterval = 50 * time.Millisecond
)

// Obstacle is a single cactus on the track.
type Obstacle struct {
	X      float64
	Height float64
	Width  float64
}

// GameState is the mutable state of one game session.
type GameState struct {
	mu        sync.RWMutex
	DinoY     float64 // CSS top offset of dino (ground = groundY)
	velY      float64
	Score     int
	Dead      bool
	Obstacles []Obstacle
	speed     float64 // obstacle scroll speed px/tick
	nextSpawn float64 // x position at which to spawn next obstacle
	clients   int     // number of active SSE watchers
}

// New creates a fresh, running GameState and starts its loop.
func New() *GameState {
	gs := &GameState{
		DinoY:     groundY,
		speed:     4.0,
		nextSpawn: arenaWidth + 200,
	}
	go gs.loop()
	return gs
}

// RegisterClient increments the watcher count.
func (gs *GameState) RegisterClient() {
	gs.mu.Lock()
	gs.clients++
	gs.mu.Unlock()
}

// UnregisterClient decrements the watcher count.
func (gs *GameState) UnregisterClient() {
	gs.mu.Lock()
	gs.clients--
	gs.mu.Unlock()
}

// Clients returns the current number of SSE clients.
func (gs *GameState) Clients() int {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.clients
}

// Jump triggers a jump if the dino is on the ground, or restarts after death.
func (gs *GameState) Jump() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if gs.Dead {
		gs.DinoY = groundY
		gs.velY = 0
		gs.Score = 0
		gs.Dead = false
		gs.Obstacles = gs.Obstacles[:0]
		gs.speed = 4.0
		gs.nextSpawn = arenaWidth + 200
		return
	}
	if gs.DinoY >= groundY {
		gs.velY = jumpVel
	}
}

// Restart resets the game after death.
func (gs *GameState) Restart() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.DinoY = groundY
	gs.velY = 0
	gs.Score = 0
	gs.Dead = false
	gs.Obstacles = gs.Obstacles[:0]
	gs.speed = 4.0
	gs.nextSpawn = arenaWidth + 200
}

// Snapshot returns a safe copy of the current state.
func (gs *GameState) Snapshot() (snap GameState, clients int) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	snap = GameState{
		DinoY:     gs.DinoY,
		Score:     gs.Score,
		Dead:      gs.Dead,
		Obstacles: make([]Obstacle, len(gs.Obstacles)),
	}
	copy(snap.Obstacles, gs.Obstacles)
	clients = gs.clients
	return
}

// loop runs the physics ticker.
func (gs *GameState) loop() {
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()
	for range ticker.C {
		gs.tick()
	}
}

func (gs *GameState) tick() {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Dead {
		return
	}

	// -- Physics: dino vertical --
	gs.velY += gravity
	gs.DinoY += gs.velY
	if gs.DinoY >= groundY {
		gs.DinoY = groundY
		gs.velY = 0
	}

	// -- Score & speed ramp --
	gs.Score++
	if gs.Score%200 == 0 && gs.speed < 14 {
		gs.speed += 0.5
	}

	// -- Obstacles: move & remove off-screen --
	alive := gs.Obstacles[:0]
	for _, o := range gs.Obstacles {
		o.X -= gs.speed
		if o.X+o.Width > -10 {
			alive = append(alive, o)
		}
	}
	gs.Obstacles = alive

	// -- Spawn new obstacle --
	rightmost := 0.0
	for _, o := range gs.Obstacles {
		if o.X > rightmost {
			rightmost = o.X
		}
	}
	if len(gs.Obstacles) < maxSlots && rightmost < gs.nextSpawn {
		h := 35 + rand.Float64()*30  // 35–65 px tall
		w := 18 + rand.Float64()*10  // 18–28 px wide
		gs.Obstacles = append(gs.Obstacles, Obstacle{
			X:      arenaWidth,
			Height: h,
			Width:  w,
		})
		// Next spawn gap: 300-600 px from current rightmost edge
		gs.nextSpawn = arenaWidth + 300 + rand.Float64()*300
	}

	// -- Collision detection --
	// Dino box: x=60..110, bottom at groundY+dinoHeight
	dinoLeft := 60.0
	dinoRight := 110.0
	dinoBottom := groundY + dinoHeight
	dinoTop := gs.DinoY

	for _, o := range gs.Obstacles {
		obsRight := o.X + o.Width
		obsTop := groundY + dinoHeight - o.Height // obstacle sits on ground
		if dinoRight > o.X && dinoLeft < obsRight &&
			dinoBottom > obsTop && dinoTop < dinoBottom {
			gs.Dead = true
			return
		}
	}
}

package handlers

// SVG canvas dimensions.
const (
	svgWidth  = 800
	svgHeight = 200
	groundY   = 160 // SVG y-coordinate of the ground line
)

// Dino geometry.
const (
	dinoX      = 80
	dinoWidth  = 28
	dinoHeight = 40
)

// Physics and gameplay tuning.
const (
	jumpImpulse = 11.0
	gravity     = 0.55
	baseSpeed   = 4.0
)

// Obstacle is a cactus on the ground.
type Obstacle struct {
	X      float64
	Width  float64
	Height float64
}

// GameState is the full mutable game state. Owned by gameMu in loop.go.
type GameState struct {
	DinoY     float64 // distance above ground (0 = standing)
	DinoVY    float64 // vertical velocity (+ve = upward)
	Obstacles []Obstacle
	Score     int
	Speed     float64
	Dead      bool
	DeadTimer int // ticks until auto-reset after death
}

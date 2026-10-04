package handlers

import "math/rand"

// applyGravity updates the dino's vertical position each tick.
func applyGravity(gs *GameState) {
	gs.DinoVY -= gravity
	gs.DinoY += gs.DinoVY
	if gs.DinoY < 0 {
		gs.DinoY = 0
		gs.DinoVY = 0
	}
}

// moveObstacles scrolls all obstacles left and removes off-screen ones.
func moveObstacles(gs *GameState) {
	for i := range gs.Obstacles {
		gs.Obstacles[i].X -= gs.Speed
	}
	var alive []Obstacle
	for _, o := range gs.Obstacles {
		if o.X+o.Width > 0 {
			alive = append(alive, o)
		}
	}
	gs.Obstacles = alive
}

// maybeSpawn adds a new obstacle once the last one is far enough left.
func maybeSpawn(gs *GameState) {
	for _, o := range gs.Obstacles {
		if o.X > float64(svgWidth)-350 {
			return
		}
	}
	gs.Obstacles = append(gs.Obstacles, Obstacle{
		X:      float64(svgWidth) + float64(rand.Intn(220)),
		Width:  float64(18 + rand.Intn(18)),
		Height: float64(30 + rand.Intn(45)),
	})
}

// checkCollision sets Dead=true if the dino overlaps any obstacle.
func checkCollision(gs *GameState) {
	dinoL := float64(dinoX) + 2
	dinoR := float64(dinoX+dinoWidth) - 2
	dinoB := float64(groundY) - gs.DinoY
	dinoT := dinoB - dinoHeight
	for _, o := range gs.Obstacles {
		if dinoR > o.X && dinoL < o.X+o.Width &&
			dinoB > float64(groundY)-o.Height && dinoT < float64(groundY) {
			gs.Dead = true
			gs.DeadTimer = 125 // ~2 s at 30 fps
			return
		}
	}
}

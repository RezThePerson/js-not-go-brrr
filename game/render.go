package game

import (
	"fmt"
	"strings"
)

const (
	groundLine = groundY + dinoHeight // y of the ground line in SVG coords
	dinoX      = 60.0
	dinoWidth  = 50.0
)

// Render returns a single SVG frame representing the current game state.
func Render(gs *GameState, playerCount int) []byte {
	var b strings.Builder
	b.Grow(1024)

	dinoColor := "#e5e7eb"
	if gs.Dead {
		dinoColor = "#ef4444"
	}

	players := fmt.Sprintf("%d player", playerCount)
	if playerCount != 1 {
		players += "s"
	}
	score := fmt.Sprintf("%05d", gs.Score)

	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 300" width="800" height="300">`)

	// Background
	b.WriteString(`<rect width="800" height="300" fill="#2c2c2c"/>`)

	// Ground
	fmt.Fprintf(&b, `<line x1="0" y1="%.1f" x2="800" y2="%.1f" stroke="#e5e7eb" stroke-width="2"/>`,
		groundLine, groundLine)

	// Dino
	fmt.Fprintf(&b,
		`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" rx="4"/>`,
		dinoX, gs.DinoY, dinoWidth, dinoHeight, dinoColor)

	// Obstacles
	for _, o := range gs.Obstacles {
		obsY := groundLine - o.Height
		fmt.Fprintf(&b,
			`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#e5e7eb" rx="2"/>`,
			o.X, obsY, o.Width, o.Height)
	}

	// Score (top-right)
	fmt.Fprintf(&b,
		`<text x="790" y="28" text-anchor="end" fill="#e5e7eb" font-family="monospace" font-size="18">%s</text>`,
		score)

	// Player count (top-left)
	fmt.Fprintf(&b,
		`<text x="10" y="28" fill="#9ca3af" font-family="monospace" font-size="14">%s</text>`,
		players)

	// Game-over overlay
	if gs.Dead {
		b.WriteString(`<rect width="800" height="300" fill="rgba(0,0,0,0.55)"/>`)
		b.WriteString(`<text x="400" y="135" text-anchor="middle" fill="#ef4444" font-family="monospace" font-size="36" font-weight="bold">GAME OVER</text>`)
		b.WriteString(`<text x="400" y="175" text-anchor="middle" fill="#e5e7eb" font-family="monospace" font-size="16">press jump to restart</text>`)
	}

	b.WriteString(`</svg>`)
	return []byte(b.String())
}

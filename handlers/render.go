package handlers

import (
	"fmt"
	"strings"
)

// Dark-mode palette baked into the SVG frames.
const (
	colorBg       = "#111827" // page/SVG background
	colorGround   = "#4b5563" // ground line
	colorDino     = "#e5e7eb" // dino body
	colorDinoDead = "#ef4444" // dino when dead
	colorObstacle = "#22c55e" // cactus
	colorText     = "#e5e7eb" // HUD text
	colorOverlay  = "#ef4444" // game-over panel
)

// RenderFrame produces an SVG image representing the current game state.
//
// SVG canvas: 800×200px. CSS on the client crops to the bottom 120px
// (showing SVG y=80..200), so all HUD elements are placed within that zone.
//
//	y=80  ┌──────────────────────────────── crop boundary ─┐
//	      │  score (y=100)         player count (y=100)    │
//	      │                                                 │
//	y=120 │  dino top when on ground                        │
//	      │  [dino body]                                    │
//	y=160 │  ─────────────── ground ──────────────────────  │
//	      │  (underground HUD strip)                        │
//	y=200 └─────────────────────────────────────────────────┘
func RenderFrame(gs *GameState, playerCount int) []byte {
	var b strings.Builder
	b.Grow(2048)

	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		svgWidth, svgHeight)

	// Background
	fmt.Fprintf(&b, `<rect width="%d" height="%d" fill="%s"/>`, svgWidth, svgHeight, colorBg)

	// Ground line
	fmt.Fprintf(&b, `<line x1="0" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2"/>`,
		groundY, svgWidth, groundY, colorGround)

	// Dino
	dinoColor := colorDino
	if gs.Dead {
		dinoColor = colorDinoDead
	}
	dinoSVGY := float64(groundY) - float64(dinoHeight) - gs.DinoY
	// Body
	fmt.Fprintf(&b, `<rect x="%d" y="%.1f" width="%d" height="%d" rx="3" fill="%s"/>`,
		dinoX, dinoSVGY, dinoWidth, dinoHeight, dinoColor)
	// Head (slightly wider, sitting on top of body)
	fmt.Fprintf(&b, `<rect x="%d" y="%.1f" width="%d" height="%d" rx="3" fill="%s"/>`,
		dinoX+4, dinoSVGY-12, dinoWidth+4, 14, dinoColor)

	// Obstacles (cacti)
	for _, o := range gs.Obstacles {
		obsY := float64(groundY) - o.Height
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="2" fill="%s"/>`,
			o.X, obsY, o.Width, o.Height, colorObstacle)
	}

	// --- HUD (all within the y=80..200 visible crop) ---

	// Score — top-left of visible area
	fmt.Fprintf(&b,
		`<text x="10" y="100" font-family="monospace" font-size="15" fill="%s">%05d</text>`,
		colorText, gs.Score)

	// Player count badge — top-right of visible area
	badge := fmt.Sprintf("%d player", playerCount)
	if playerCount != 1 {
		badge += "s"
	}
	fmt.Fprintf(&b,
		`<text x="790" y="100" font-family="monospace" font-size="15" fill="%s" text-anchor="end">&#x1F465; %s</text>`,
		colorText, badge)

	// Game-over overlay
	if gs.Dead {
		b.WriteString(`<rect x="250" y="110" width="300" height="50" rx="6" fill="#ef4444" opacity="0.9"/>`)
		b.WriteString(`<text x="400" y="142" font-family="monospace" font-size="22" fill="white" text-anchor="middle" font-weight="bold">GAME OVER</text>`)
	}

	b.WriteString(`</svg>`)
	return []byte(b.String())
}

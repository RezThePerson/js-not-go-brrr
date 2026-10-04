package handlers

import (
	"fmt"
	"strings"
)

// RenderFrame renders a complete SVG frame for the current game state.
//
// The SVG is 800×200px. The client CSS crops it to the bottom 120px
// (y=80..200), keeping the ground, dino, and HUD always in view.
func RenderFrame(gs *GameState, playerCount int) []byte {
	var b strings.Builder
	b.Grow(2048)

	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`, svgWidth, svgHeight)
	fmt.Fprintf(&b, `<rect width="%d" height="%d" fill="%s"/>`, svgWidth, svgHeight, colorBg)
	fmt.Fprintf(&b, `<line x1="0" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="2"/>`,
		groundY, svgWidth, groundY, colorGround)

	drawDino(&b, gs)
	drawObstacles(&b, gs)
	drawHUD(&b, gs, playerCount)

	b.WriteString(`</svg>`)
	return []byte(b.String())
}

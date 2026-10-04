package handlers

import (
	"fmt"
	"strings"
)

// drawDino writes dino body + head rectangles.
func drawDino(b *strings.Builder, gs *GameState) {
	col := colorDino
	if gs.Dead {
		col = colorDinoDead
	}
	y := float64(groundY) - float64(dinoHeight) - gs.DinoY
	fmt.Fprintf(b, `<rect x="%d" y="%.1f" width="%d" height="%d" rx="3" fill="%s"/>`,
		dinoX, y, dinoWidth, dinoHeight, col)
	fmt.Fprintf(b, `<rect x="%d" y="%.1f" width="%d" height="%d" rx="3" fill="%s"/>`,
		dinoX+4, y-12, dinoWidth+4, 14, col)
}

// drawObstacles writes all cactus rectangles.
func drawObstacles(b *strings.Builder, gs *GameState) {
	for _, o := range gs.Obstacles {
		fmt.Fprintf(b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="2" fill="%s"/>`,
			o.X, float64(groundY)-o.Height, o.Width, o.Height, colorObstacle)
	}
}

// drawHUD writes score, player-count badge, and game-over overlay.
func drawHUD(b *strings.Builder, gs *GameState, players int) {
	fmt.Fprintf(b, `<text x="10" y="100" font-family="monospace" font-size="15" fill="%s">%05d</text>`,
		colorText, gs.Score)
	badge := fmt.Sprintf("%d player", players)
	if players != 1 {
		badge += "s"
	}
	fmt.Fprintf(b, `<text x="790" y="100" font-family="monospace" font-size="15" fill="%s" text-anchor="end">&#x1F465; %s</text>`,
		colorText, badge)
	if gs.Dead {
		fmt.Fprintf(b, `<rect x="250" y="110" width="300" height="50" rx="6" fill="%s" opacity="0.9"/>`, colorOverlay)
		b.WriteString(`<text x="400" y="142" font-family="monospace" font-size="22" fill="white" text-anchor="middle" font-weight="bold">GAME OVER</text>`)
	}
}

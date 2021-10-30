package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

const (
	// Update is still ~60 (default) TPS to listen better for mouse events, this just applies to game logic
	gameUpdateOnFrame = 20
)

var (
	gameFrameCount   = 0
	bgColor, _       = colorx.ParseHexColor("#343a40")
	bgColorPaused, _ = colorx.ParseHexColor("#343a50")
)

type Game struct {
	grid   *Grid
	paused bool
}

func (g *Game) BgColor() color.RGBA {
	if g.paused {
		return bgColorPaused
	}

	return bgColor
}

func (g *Game) Update() error {
	inputUpdate()

	// Only update game on every 20 frames
	gameFrameCount = (gameFrameCount + 1) % gameUpdateOnFrame
	if gameFrameCount == 0 {
		gameUpdate()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen, g.BgColor())
	drawDots(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g Game) Paused() bool {
	return g.paused
}

func (g *Game) Pause() {
	g.paused = true
}

func (g *Game) UnPause() {
	g.paused = false
}

func (g *Game) TogglePause() {
	g.paused = !g.paused
}

func (g *Game) Restart() {
	g.grid.Clear()
}

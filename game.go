package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

const (
	// Update is still ~60 (default) TPS to listen better for mouse events, this just applies to game logic
	GAME_TICK_FRAME_NUM = 15
)

var (
	gameFrameCount   = 0
	bgColor, _       = colorx.ParseHexColor("#303040")
	bgColorPaused, _ = colorx.ParseHexColor("#3f3f4a")
	bgCellColor, _   = colorx.ParseHexColor("#022330")
)

type Game struct {
	grid       *Grid
	paused     bool
	generation int
}

func (g Game) BgColor() color.RGBA {
	if g.paused {
		return bgColorPaused
	}

	return bgColor
}

func (g Game) BgCellColor() color.RGBA {
	return bgCellColor
}

func (g *Game) Update() error {
	inputUpdate()

	// Only update game on every 20 frames
	gameFrameCount = (gameFrameCount + 1) % GAME_TICK_FRAME_NUM
	if gameFrameCount == 0 {
		gameUpdate()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen, g.BgColor())
	drawDots(screen)
	drawOverlay(screen, g.BgCellColor())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
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
	g.generation = 0
	g.grid.Clear()
}

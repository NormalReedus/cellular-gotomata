package main

import (
	"image/color"
	_ "image/png" // necessary for loading images
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

//TODO: figure out how to easily define a pixel size that everything multiplies by, so we can have cells, that are actually many pixels (which allows us to draw borders on every cell, so you can see the grid)

const (
	screenWidth, screenHeight = 30, 30
	pixelSize                 = 30
)

var (
	game *Game
)

func init() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth*pixelSize, screenHeight*pixelSize)
	ebiten.SetWindowTitle("Cellular Automata")
}

func setupInitialState() {
	game = &Game{grid: NewGrid(), paused: true}
}

func main() {
	setupInitialState()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Default TPS
func inputUpdate() {
	//TODO: make these handlers into functions
	coords := leftClick()
	if coords != nil {
		NewDot(*coords, game.grid)
	}

	coords = rightClick()
	if coords != nil {
		dot, err := game.grid.Get(*coords)
		if dot == nil || err != nil {
			return
		}

		dot.Remove()
	}

	if spaceKey() {
		game.TogglePause()
	}

	if cKey() {
		game.Restart()
	}
}

// Slower TPS
func gameUpdate() {
	if game.Paused() {
		return
	}

	gol := NewMyGameOfLife()
	game.grid.Convolve(gol)

	game.generation++
}

func drawBackground(screen *ebiten.Image, clr color.RGBA) {
	screen.Fill(clr)
	// ebitenutil.DebugPrint(screen, fmt.Sprint(game.generation)) //TODO: Do this again when text can be smaller
}

func drawDots(screen *ebiten.Image) {
	game.grid.ForEach(func(dot *Dot) {
		dot.Draw(screen)
	})
}

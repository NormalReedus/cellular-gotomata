package main

import (
	"fmt"
	"image/color"
	_ "image/png" // necessary for loading images
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//TODO: figure out how to easily define a pixel size that everything multiplies by, so we can have cells, that are actually many pixels (which allows us to draw borders on every cell, so you can see the grid)

const (
	GRID_WIDTH, GRID_HEIGHT     = 60, 60 //* BOTH MUST BE DIVISIBLE BY CELL_SIZE
	CELL_SIZE                   = 20
	SCREEN_WIDTH, SCREEN_HEIGHT = GRID_WIDTH * CELL_SIZE, GRID_HEIGHT * CELL_SIZE
)

var (
	game *Game
	gol  Convolver
	// golMod Convolver
)

func init() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Cellular Automata")
}

func setupInitialState() {
	game = &Game{grid: NewGrid(), paused: true}
	gol = NewCustomGame2()
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

	// Doing more convolutions per tick can 'modify' an existing Game of Life to compose brand new games
	game.grid.Convolve(gol)

	game.generation++
}

func drawBackground(screen *ebiten.Image, clr color.RGBA) {
	screen.Fill(clr)
}

func drawDots(screen *ebiten.Image) {
	game.grid.ForEach(func(dot *Dot) {
		dot.Draw(screen)
	})
}

func drawOverlay(screen *ebiten.Image, bgCellColor color.RGBA) {

	for x := 0; x < SCREEN_WIDTH; x += CELL_SIZE {
		ebitenutil.DrawLine(screen, float64(x), 0, float64(x), float64(SCREEN_HEIGHT), bgCellColor)
	}
	for y := 0; y < SCREEN_HEIGHT; y += CELL_SIZE {
		ebitenutil.DrawLine(screen, 0, float64(y), float64(SCREEN_HEIGHT), float64(y), bgCellColor)
	}

	// Print generation num
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Generation: %d", game.generation))
}

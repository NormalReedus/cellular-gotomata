package main

import (
	"image/color"
	_ "image/png" // necessary for loading images
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

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

	game.grid.Convolve(3, func(win *Window) *Dot {
		//* https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life
		dead := win.NumEmptyNeighbors()
		alive := 8 - dead

		currentVal := win.Center()

		// "Any live cell with two or three live neighbours survives."
		if currentVal != nil && between(alive, 2, 3) {
			return currentVal
		}

		// "Any dead cell with three live neighbours becomes a live cell."
		if currentVal == nil && alive == 3 {
			return NewDot(win.GridCoords(), nil) // parentGrid is set in grid.Convolve instead
		}

		// "All other live cells die in the next generation." / "...all other dead cells stay dead."
		return nil
	})

}

func drawBackground(screen *ebiten.Image, clr color.RGBA) {
	screen.Fill(clr)
}

func drawDots(screen *ebiten.Image) {
	game.grid.ForEach(func(dot *Dot) {
		dot.Draw(screen)
	})
}

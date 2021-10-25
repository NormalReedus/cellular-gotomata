package main

import (
	"image/color"
	_ "image/png" // necessary for loading images
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

//TODO: In GameUpdate - Convolve through grid, implement GoL in callback, replace game.grid with resulting tempGrid

const (
	screenWidth, screenHeight = 30, 30
	numDots                   = 3
	// Update is still ~60 (default) TPS to listen better for mouse events, this just applies to game logic
	gameUpdateOnFrame = 20
)

var (
	gameFrameCount   = 0
	bgColor, _       = colorx.ParseHexColor("#343a40")
	bgColorPaused, _ = colorx.ParseHexColor("#343a50")
	game             *Game
)

func init() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth*30, screenHeight*30)
	ebiten.SetWindowTitle("Cellular Automata")

	debugInit()
}

func setupInitialState() {
	game = &Game{grid: NewGrid(), list: NewLinkedList(), paused: true}
}

func main() {
	setupInitialState()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	grid   *Grid
	list   *LinkedList
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

	debugUpdate()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.BgColor())

	drawDots(screen)

	debugDraw()
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

// Default TPS
func inputUpdate() {
	//TODO: make these handlers into functions
	coords := leftClick()
	if coords != nil {
		NewDot(*coords, game.list, game.grid)
	}

	coords = rightClick()
	if coords != nil {
		node := game.grid.Get(*coords)
		if node == nil {
			return
		}
		dot := node.(*Dot)

		dot.Remove()
	}

	if spaceKey() {
		game.TogglePause()
	}
}

// Slower TPS
func gameUpdate() {
	if game.Paused() {
		return
	}

}

func drawDots(screen *ebiten.Image) {
	game.list.ForEach(func(nm NodeManipulator) {
		// Assert that dotnode is a *Dot, so we can use *Dot's methods
		dot := nm.(*Dot)
		dot.Draw(screen)
	}, false)
}

// func wanderDots() {
// 	game.mainList.ForEach(func(nm NodeManipulator) {
// 		// Assert that dotnode is a *Dot, so we can use *Dot's methods
// 		dot := nm.(*Dot)
// 		newCoords, _ := game.mainGrid.RandomOpenCell()

// 		if err := dot.MoveCell(*newCoords); err != nil {
// 			fmt.Println(err) // should never happen, since RandomOpenCell should never return an occupied cell
// 		}
// 	}, false)
// }

func debugInit() {

}

func debugUpdate() {

}

func debugDraw() {
}

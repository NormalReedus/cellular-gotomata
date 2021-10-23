package main

import (
	_ "image/png" // necessary for loading images
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

//TODO: there should be something that adds a dot to both grid and ll at the same time (a method on Game? Game would then also have a prop for grid and ll that is THE structures to add to). It should take coords, create the Dot, call ll.Add and grid.Set with the coords
//TODO: add struct / interfaces for rules of the game. There should be different structs that implement the same interface, that specify the rules for every tick. This should be embedded in Dot and Grid and LL etc.
// There could for example be a Move() in the interface, and depending on the struct, the Move() function behaves differently, e.g. by moving a specific number of steps specified by a struct field or whether it should die when something specific happens etc.
const (
	screenWidth, screenHeight = 30, 30
	numDots                   = 3
	// Update is still ~60 (default) TPS to listen better for mouse events, this just applies to game logic
	gameUpdateOnFrame = 20
)

var (
	gameFrameCount = 0
	bgColor, _     = colorx.ParseHexColor("#343a40")
	game           *Game
)

func init() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth*30, screenHeight*30)
	ebiten.SetWindowTitle("Cellular Automata")

	debugInit()
}

func setupInitialState() {
	game = &Game{dotList: NewLinkedList(), dotGrid: NewGrid()}

	for i := 0; i < numDots; i++ {
		coords, err := game.dotGrid.RandomOpenCell()
		if err != nil {
			log.Fatal("you have added too many dots for this grid")
		}

		NewDot(*coords, game.dotList, game.dotGrid)
	}
}

func main() {
	setupInitialState()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	dotList *LinkedList
	dotGrid *Grid
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

// Default TPS
func inputUpdate() {
	//TODO: make these handlers into functions
	coords := leftClick()
	if coords != nil {
		NewDot(*coords, game.dotList, game.dotGrid)

	}
	coords = rightClick()
	if coords != nil {
		node := game.dotGrid.Get(*coords)
		if node == nil {
			return
		}
		dot := node.(*Dot)

		dot.Remove()
	}
}

// Slower TPS
func gameUpdate() {
	// wanderDots()
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)

	drawDots(screen)

	debugDraw()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	return screenWidth, screenHeight
}

func drawDots(screen *ebiten.Image) {
	game.dotList.ForEach(func(nm NodeManipulator) {
		// Assert that dotnode is a *Dot, so we can use *Dot's methods
		dot := nm.(*Dot)
		dot.Draw(screen)
	}, false)
}

func wanderDots() {
	game.dotList.ForEach(func(nm NodeManipulator) {
		// Assert that dotnode is a *Dot, so we can use *Dot's methods
		dot := nm.(*Dot)
		newCoords, _ := game.dotGrid.RandomOpenCell()
		dot.MoveCell(*newCoords)

	}, false)
}

func debugInit() {
}

func debugUpdate() {

}

func debugDraw() {
}

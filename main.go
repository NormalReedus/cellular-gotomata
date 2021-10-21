package main

import (
	"fmt"
	_ "image/png" // necessary for loading images
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/icza/gox/imagex/colorx"
)

//TODO: LOG OUT DOTS, LINKEDLIST to see if they print as they are instructed and all fields work and test that everything is being passed as pointer / value correctly when specifying interfaces
//TODO: continue with Grid and Cell the same as with LinkedList and Node, embedding Cell into Dot to implement the CellManipulator interface that together with Cell have the fields and methods to interact with the Grid
//TODO: grid should then use CellManipulator for everything, since Dot should implement CellManipulator. Grid should just be another way to reference Dots with O(1) lookups

type Game struct{}

const (
	screenWidth, screenHeight = 30, 30
	numDots                   = 10
	// Update is still ~60 (default) TPS to listen better for mouse events, this just applies to game logic
	gameUpdateOnFrame = 20
)

var (
	gameFrameCount = 0
	bgColor, _     = colorx.ParseHexColor("#343a40")
	dots           *LinkedList
)

func init() {
	var dotSlice []NodeManipulator

	for i := 0; i < numDots; i++ {
		dot := NewDot(screenWidth-1, screenHeight-1, screenWidth, screenHeight)
		dotSlice = append(dotSlice, dot)
	}

	dots = NewLinkedList(dotSlice...)

	debugInit()
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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		fmt.Println(ebiten.CursorPosition())
	}
}

// Slower TPS
func gameUpdate() {
	wanderDots()
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)

	drawDots(screen)

	debugDraw()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{}

	ebiten.SetWindowSize(screenWidth*30, screenHeight*30)
	ebiten.SetWindowTitle("Cellular Automata")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func drawDots(screen *ebiten.Image) {
	dots.ForEach(func(nm NodeManipulator) {
		// Assert that dotnode is a *Dot, so we can use *Dot's methods
		dot := nm.(*Dot)
		dot.Draw(screen)
	}, false)
}

func wanderDots() {
	dots.ForEach(func(nm NodeManipulator) {
		// Assert that dotnode is a *Dot, so we can use *Dot's methods
		dot := nm.(*Dot)
		// must be -1 to not go outside window
		x := float64(rand.Intn(screenWidth - 1))
		y := float64(rand.Intn(screenHeight - 1))

		dot.Position.Set(x, y)
	}, false)
}

func debugInit() {
}

func debugUpdate() {

}

func debugDraw() {
}

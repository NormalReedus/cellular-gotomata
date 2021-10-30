package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

//* -------------------------
//* GRID
//* -------------------------
type Grid struct {
	data         ScreenPixelMatrix
	numUsedCells int
}

func NewGrid() *Grid {
	return &Grid{}
}

func (g *Grid) String() string {
	return fmt.Sprintf("Grid{ numUsedCells: %d }", g.numUsedCells)
}

func (g *Grid) Clear() {
	g.ReplaceMatrix(g.CreateTempMatrix())
	g.numUsedCells = 0
}

func (g Grid) CreateTempMatrix() ScreenPixelMatrix {
	var matrix ScreenPixelMatrix

	return matrix
}

// Returns the last col and row num that can contain values
func (g *Grid) Bounds() (int, int) {
	return len(g.data) - 1, len(g.data[0]) - 1
}

//TODO: return err (or nil?) if the coords are out of bounds
func (g *Grid) Get(coords Point) (*Dot, error) {
	maxX, maxY := g.Bounds()
	if !between(coords.X, 0, maxX) || !between(coords.Y, 0, maxY) {
		return nil, fmt.Errorf("cannot 'Get' cell out of grid bounds. Accessing: %v of x: 0 - %d, y: 0 - %d", coords, maxX, maxY)
	}

	return g.data[coords.X][coords.Y], nil
}

func (g *Grid) Move(currentCoords Point, newCoords Point, dot *Dot) error {
	targetDot, err := g.Get(newCoords)

	// Cell already has a dot
	if targetDot != nil {
		currentDot, _ := g.Get(currentCoords)
		return fmt.Errorf("cannot move %v to cell with coords %v: the target cell is already occupied", currentDot, newCoords)
	}
	// Tried to access out of bounds
	if err != nil {
		currentDot, _ := g.Get(currentCoords)
		return fmt.Errorf("cannot move %v to cell with coords %v: the target cell is already occupied", currentDot, newCoords)
	}

	g.Remove(currentCoords)
	g.Set(newCoords, dot)

	return nil
}

// Not really used with tempMatrix (in convolutions)
func (g *Grid) Set(coords Point, dot *Dot) {
	// Remove existing (if any) dot first, to also decrement usedCells etc
	if dot, _ := g.Get(coords); dot != nil {
		g.Remove(coords)
	}

	g.data[coords.X][coords.Y] = dot

	dot.SetParentGrid(g)
	dot.SetPosition(coords)

	g.IncrementNumUsedCells()
}

func (g *Grid) ReplaceMatrix(data ScreenPixelMatrix) {
	g.data = data
}

func (g *Grid) Remove(coords Point) {
	if dot, _ := g.Get(coords); dot == nil {
		return
	}

	g.data[coords.X][coords.Y] = nil

	g.DecrementNumUsedCells()
}

func (g *Grid) RandomOpenCell() (*Point, error) {
	var openCells []Point

	for x, column := range g.data {
		for y := range column {
			if g.data[x][y] == nil {
				p := Point{X: x, Y: y}
				openCells = append(openCells, p)
			}
		}
	}

	if len(openCells) == 0 {
		return nil, errors.New("there are no more open cells")
	}

	randCellNum := rand.Intn(len(openCells))

	return &openCells[randCellNum], nil
}

func (g *Grid) IncrementNumUsedCells() {
	g.numUsedCells++
}

func (g *Grid) DecrementNumUsedCells() {
	g.numUsedCells--
}

// Loop though all cells in grid and do an operation within a window (e.g. a kernel operation)
func (g *Grid) Convolve(windowSize int, callback func(*Window) *Dot) {
	tempMatrix := g.CreateTempMatrix()

	for x := 0; x < len(g.data); x++ {
		for y := 0; y < len(g.data[x]); y++ {
			coords := *NewPoint(x, y)

			win := NewWindow(g, coords, windowSize)

			cellVal := callback(win)

			if cellVal != nil {
				cellVal.SetPosition(coords)
			}

			tempMatrix[x][y] = cellVal

			// Keep track of removed / added cells
			formerCellVal, _ := g.Get(coords)
			// If this cell used to have a value, but now doesn't
			if formerCellVal != nil && cellVal == nil {
				g.DecrementNumUsedCells()
			}
			// If this cell had no value, but now it does
			if formerCellVal == nil && cellVal != nil {
				g.IncrementNumUsedCells()
			}
		}
	}

	g.ReplaceMatrix(tempMatrix)

	// New dots made in the callback should not have a parentGrid and position in that grid yet, since it is born into the tempMatrix instead
	// ...therefore we need to set the parent grid for every (newly created) dot here
	// Other g.Set() functionality such as decrement/increment usedCells and setting position is handled separately above
	// ... as to avoid using g.Set()
	g.ForEach(func(dot *Dot) {
		dot.parentGrid = g
	})
}

//*NOTE: collisions can happen, if no intermediary temp matrix is used
func (g *Grid) ForEach(callback func(dot *Dot)) {
	// Adding all existing dots to a slice up front makes sure we will only call callback on every dot once
	var dots []*Dot

	for x := 0; x < len(g.data); x++ {
		for y := 0; y < len(g.data[x]); y++ {
			if g.data[x][y] == nil {
				continue
			}

			dots = append(dots, g.data[x][y])
		}
	}

	if len(dots) == 0 {
		return
	}

	for _, dot := range dots {
		callback(dot)
	}
}

//* -------------------------
//* WINDOW
//* -------------------------
type Window struct {
	grid   *Grid
	center Point
	size   int
	matrix [][]*Dot
}

// Will pad the grid with nil values if needed
func NewWindow(grid *Grid, coords Point, size int) *Window {
	if size%2 != 1 {
		log.Fatal("window 'size' can only be an odd number")
	}

	window := &Window{grid: grid, center: coords, size: size}

	var matrix [][]*Dot // matches ScreenPixelMatrix's type, but this is dynamically sized to window size instead of array

	reach := window.Reach()
	winMinX, winMaxX := coords.X-reach, coords.X+reach
	winMinY, winMaxY := coords.Y-reach, coords.Y+reach

	for x := winMinX; x <= winMaxX; x++ {
		col := make([]*Dot, 0)

		for y := winMinY; y <= winMaxY; y++ {
			// Can be nil
			cellValue, _ := grid.Get(*NewPoint(x, y))

			col = append(col, cellValue)
		}

		matrix = append(matrix, col)
	}

	window.matrix = matrix

	return window
}

// Returns how many cells to each side of center the window spans
func (w Window) Reach() int {
	return w.size >> 1 // divide by 2, round down
}

// Returns the index number of the center cell inside window (x AND y coord of center)
func (w Window) CenterIndex() int {
	return w.size - w.Reach() - 1
}

// Returns the window's centerpoint's value
func (w *Window) Center() *Dot {
	index := w.CenterIndex()
	return w.matrix[index][index]
}

// Returns value of arbitrary coord inside window
func (w *Window) Get(coords Point) *Dot {
	return w.matrix[coords.X][coords.Y]
}

// Returns the number of empty cells around the center cell
func (w Window) NumEmptyNeighbors() int {

	var count int

	center := w.CenterIndex()

	for x, col := range w.matrix {
		for y, val := range col {
			if x == center && y == center {
				continue
			}

			if val == nil {
				count++
			}
		}
	}

	return count
}

func (w *Window) AliveNeighbors() []*Dot {
	var dots []*Dot

	center := w.CenterIndex()

	for x, col := range w.matrix {
		for y, val := range col {
			if x == center && y == center {
				continue
			}

			if val != nil {
				dots = append(dots, val)
			}
		}
	}

	return dots
}

// Returns the grid coords of the center cell of the window
func (w *Window) GridCoords() Point {
	return w.center
}

//* -------------------------
//* DOT
//* -------------------------
type Dot struct {
	image      *ebiten.Image
	fill       color.Color
	parentGrid *Grid
	position   Point
}

// Set parentGrid to nil to not immediately add to a grid (in convolutions etc)
func NewDot(coords Point, parentGrid *Grid) *Dot {
	image := ebiten.NewImage(1, 1)
	color, err := colorx.ParseHexColor("#adb5bd")
	if err != nil {
		log.Fatal(err)
	}

	dot := &Dot{
		image:    image,
		position: coords,
		fill:     color,
	}

	if parentGrid != nil {
		dot.SetParentGrid(parentGrid)
		parentGrid.Set(coords, dot)
	}

	return dot
}

// Pretty print the Dot position x & y coordinates
func (d Dot) String() string {
	return fmt.Sprintf("Dot{ x: %d, y: %d }", d.Position().X, d.Position().Y)
}

// Grid will handle updating position of Dot if everything goes well
func (d *Dot) MoveTo(coords Point) error {
	return d.parentGrid.Move(d.position, coords, d)
}

func (d *Dot) Remove() {
	d.parentGrid.Remove(d.position)
}

func (d *Dot) SetParentGrid(grid *Grid) {
	d.parentGrid = grid
}

func (d *Dot) Position() Point {
	return d.position
}

func (d *Dot) SetPosition(coords Point) {
	d.position.SetCoords(coords.X, coords.Y)
}

func (d *Dot) Draw(screen *ebiten.Image) {
	dotOpts := &ebiten.DrawImageOptions{}
	dotOpts.GeoM.Translate(float64(d.Position().X), float64(d.Position().Y)) // position
	d.image.Fill(d.fill)                                                     // color
	screen.DrawImage(d.image, dotOpts)
}

//* -------------------------
//* POINT
//* -------------------------
type Point struct {
	X, Y int
}

func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y}
}

func (p Point) String() string {
	return fmt.Sprintf("{ x: %d, y: %d }", p.X, p.Y)
}

func (p *Point) SetCoords(x, y int) {
	p.X = x
	p.Y = y
}

func (p Point) Coords() (int, int) {
	return p.X, p.Y
}

//* -------------------------
//* SCREEN PIXEL MATRIX
//* -------------------------
type ScreenPixelMatrix [screenWidth][screenHeight]*Dot

func (spm *ScreenPixelMatrix) GetAllNonEmpty() []*Dot {
	var cells []*Dot

	for _, col := range spm {
		for _, cell := range col {
			// cell should be pointer to CellManipulator value struct
			if cell != nil {
				cells = append(cells, cell)
			}
		}
	}

	return cells
}

package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"reflect"

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

func (g Grid) CreateTempGrid() ScreenPixelMatrix {
	var tempGrid ScreenPixelMatrix

	return tempGrid
}

// Returns the last cells that can contain values
func (g *Grid) Bounds() (int, int) {
	return len(g.data) - 1, len(g.data[0]) - 1
}

func (g *Grid) Get(coords Point) CellManipulator {
	return g.data[coords.X][coords.Y]
}

func (g *Grid) Move(currentCoords Point, newCoords Point, cell CellManipulator) error {
	if g.Get(newCoords) != nil {
		return fmt.Errorf("cannot move %v to cell with coords %v: the target cell is already occupied", g.Get(currentCoords), newCoords)
	}

	g.Remove(currentCoords)
	g.Set(newCoords, cell)

	cell.SetPosition(newCoords)

	return nil
}

func (g *Grid) Set(coords Point, cell CellManipulator) {
	g.data[coords.X][coords.Y] = cell

	cell.SetParentGrid(g)

	g.IncrementNumUsedCells()
}

func (g *Grid) ReplaceState(data ScreenPixelMatrix) {
	g.data = data
}

// Only used to clear Dot in grid, to completely delete Dot, use Dot.Remove()
func (g *Grid) Remove(coords Point) {
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

func (g *Grid) Convolve(windowSize int, callback func(*Window) CellManipulator) ScreenPixelMatrix {
	tempGrid := g.CreateTempGrid()

	for x := 0; x < len(g.data); x++ {
		for y := 0; y < len(g.data[x]); y++ {
			coords := *NewPoint(x, y)

			win := NewWindow(g, coords, windowSize)

			cellVal := callback(win)
			if cellVal != nil {
				cellVal.SetPosition(coords)
			}

			tempGrid[x][y] = cellVal
		}
	}

	return tempGrid
}

//* -------------------------
//* WINDOW
//* -------------------------
type Window struct {
	grid   *Grid
	center Point
	size   int
	data   [][]CellManipulator
}

// Will pad the grid with nil values if needed
func NewWindow(grid *Grid, coords Point, size int) *Window {
	if size%2 != 1 {
		log.Fatal("window 'size' can only be an odd number")
	}

	window := &Window{grid: grid, center: coords, size: size}

	var data [][]CellManipulator

	reach := window.Reach()
	winMinX, winMaxX := coords.X-reach, coords.X+reach
	winMinY, winMaxY := coords.Y-reach, coords.Y+reach
	boundsX, boundsY := grid.Bounds()

	for x := winMinX; x <= winMaxX; x++ {
		col := make([]CellManipulator, 0)

		for y := winMinY; y <= winMaxY; y++ {

			var cellValue CellManipulator

			if x < 0 || x > boundsX || y < 0 || y > boundsY {
				cellValue = nil
			} else {
				cellValue = grid.Get(*NewPoint(x, y))
			}

			col = append(col, cellValue)
		}

		data = append(data, col)
	}

	window.data = data

	return window
}

// Returns how many cells to each side of center the window spans
func (w *Window) Reach() int {
	return w.size >> 1 // divide by 2, round down
}

func (w *Window) CenterIndex() int {
	return w.size - w.Reach() - 1
}

func (w *Window) Center() CellManipulator {
	index := w.CenterIndex()
	return w.data[index][index]
}

func (w *Window) Get(coords Point) CellManipulator {
	return w.data[coords.X][coords.Y]
}

// Returns the number of empty cells around the center cell
func (w *Window) NumEmpty() int {
	var count int

	center := w.CenterIndex()

	for x, col := range w.data {
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

// // Applying a kernel on the window should return a new value for the center position of the window
// func (w *Window) ApplyKernel() CellManipulator {
// 	return
// }

//* -------------------------
//* CELL
//* -------------------------
// Abstract struct that is embedded into Dot (i.e. not used directly anywhere)
// This makes any embedding struct implement the CellManipulator
type CellManipulator interface {
	SetParentGrid(grid *Grid)
	MoveCell(coords Point) error
	RemoveCell()
	Position() *Point
	SetPosition(Point)
}

type Cell struct {
	parentGrid *Grid
	position   Point
}

func (c *Cell) SetParentGrid(grid *Grid) {
	c.parentGrid = grid
}

func (c *Cell) RemoveCell() {
	c.parentGrid.Remove(c.position)
}

func (c *Cell) Position() *Point {
	return &c.position
}

func (c *Cell) SetPosition(coords Point) {
	c.position.SetCoords(coords.X, coords.Y)
}

//* -------------------------
//* LINKED LIST
//* -------------------------

type LinkedList struct {
	head   NodeManipulator
	tail   NodeManipulator
	length int
}

func NewLinkedList() *LinkedList {
	ll := &LinkedList{}

	return ll
}

func NewLinkedListFromMatrix(matrix *ScreenPixelMatrix) *LinkedList {
	ll := &LinkedList{}

	ll = matrix.ExportToLinkedList(ll)

	return ll
}

func (ll LinkedList) String() string {
	if ll.head != nil && ll.tail != nil {
		return fmt.Sprintf("LinkedList{ nodeType: %s, length: %d, head: %v, tail: %v }", reflect.TypeOf(ll.head), ll.length, ll.head, ll.tail)
	}

	return "LinkedList{ empty }"
}

func (ll *LinkedList) Head() NodeManipulator {
	return ll.head
}

func (ll *LinkedList) SetHead(node NodeManipulator) {
	ll.head = node
}

func (ll *LinkedList) Tail() NodeManipulator {
	return ll.tail
}

func (ll *LinkedList) SetTail(node NodeManipulator) {
	ll.tail = node

}

func (ll *LinkedList) Length() int {
	return ll.length
}

func (ll *LinkedList) Add(node NodeManipulator) {
	if ll.head == nil {
		ll.head = node
		ll.tail = node
		ll.incrementLength()

		node.SetParentList(ll)

		return
	}

	node.SetPrevNode(ll.tail)

	ll.tail.SetNextNode(node)
	ll.tail = node
	ll.incrementLength()

	node.SetParentList(ll)
}

func (ll *LinkedList) incrementLength() {
	ll.length++
}

func (ll *LinkedList) decrementLength() {
	ll.length--
}

func (ll LinkedList) ForEach(callback func(node NodeManipulator), reverse bool) {
	if reverse {
		node := ll.tail

		// Don't do anything if there are no nodes
		if node == nil {
			return
		}

		for {
			callback(node)

			if node.PrevNode() == nil {
				break
			}

			node = node.PrevNode()
		}

	} else {
		node := ll.head

		// Don't do anything if there are no nodes
		if node == nil {
			return
		}

		for {
			callback(node)

			if node.NextNode() == nil {
				break
			}

			node = node.NextNode()
		}

	}
}

//* -------------------------
//* NODE
//* -------------------------

type NodeManipulator interface {
	SetParentList(list *LinkedList)
	PrevNode() NodeManipulator
	SetPrevNode(NodeManipulator)
	NextNode() NodeManipulator
	SetNextNode(NodeManipulator)
	RemoveNode()
}

// Abstract struct that is embedded into Dot (i.e. not used directly anywhere)
// This makes any embedding struct implement the NodeManipulator
type Node struct {
	parentList *LinkedList
	prev       NodeManipulator
	next       NodeManipulator
}

func (n *Node) String() string {
	// return fmt.Sprintf("Node{ parentList: %v, prev: %v, next: %v }", n.parentList, n.prev, n.next)
	return fmt.Sprintf("Node{ prev: %v, next: %v }", n.prev, n.next)
}

func (n *Node) SetParentList(list *LinkedList) {
	n.parentList = list
}

func (n *Node) PrevNode() NodeManipulator {
	return n.prev
}

func (n *Node) SetPrevNode(node NodeManipulator) {
	n.prev = node
}

func (n *Node) NextNode() NodeManipulator {
	return n.next
}

func (n *Node) SetNextNode(node NodeManipulator) {
	n.next = node
}

func (n *Node) RemoveNode() {
	// There are always 2 refs to delete to garbage collect this node...
	if n.prev == nil {
		//* If this node is head AND tail
		if n.next == nil {
			// Both refs are from list (head, tail), since there are no other nodes
			n.parentList.SetTail(nil)
			n.parentList.SetHead(nil)

			n.parentList.decrementLength()
			return
		}

		//* If this node is ONLY head
		// One ref from list (head) and one ref from next node (prev)
		n.parentList.SetHead(n.next)
		n.next.SetPrevNode(nil)

		n.parentList.decrementLength()
		return
	}

	//* If this node is ONLY tail
	if n.next == nil {
		// One ref from list (tail) and one ref from prev node (next)
		n.parentList.SetTail(n.prev)
		n.prev.SetNextNode(nil)

		n.parentList.decrementLength()
		return
	}

	//* If this node is NEITHER head nor tail
	// One ref from both prev (next) and next (prev)
	n.prev.SetNextNode(n.next)
	n.next.SetPrevNode(n.prev)

	n.parentList.decrementLength()
}

//* -------------------------
//* DOT
//* -------------------------

type Dot struct {
	image *ebiten.Image
	fill  color.Color
	Node
	Cell
}

func NewDot(coords Point, parentList *LinkedList, parentGrid *Grid) *Dot {
	image := ebiten.NewImage(1, 1)
	color, err := colorx.ParseHexColor("#adb5bd")
	if err != nil {
		log.Fatal(err)
	}

	dot := &Dot{
		image: image,
		Cell:  Cell{position: coords},
		fill:  color,
	}

	// Remove existing dot at same location, if there is any
	existingCell := parentGrid.Get(coords)
	if existingCell != nil {
		existingDot := existingCell.(*Dot)

		existingDot.Remove()
	}

	parentList.Add(dot)
	parentGrid.Set(coords, dot)

	return dot
}

// Pretty print the Dot position x & y coordinates
func (d Dot) String() string {
	var prevPos, nextPos *Point
	if d.prev != nil {
		prevPos = d.prev.(*Dot).Position()
	}
	if d.next != nil {
		nextPos = d.next.(*Dot).Position()
	}
	return fmt.Sprintf("Dot{ x: %d, y: %d, prev: %v, next: %v }", d.Position().X, d.Position().Y, prevPos, nextPos)
}

func (d *Dot) MoveCell(coords Point) error {
	return d.parentGrid.Move(d.position, coords, d)
}

func (d *Dot) Remove() {
	d.RemoveNode()
	d.RemoveCell()
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
// Abstract struct that is embedded into Dot (i.e. not used directly anywhere)
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
type ScreenPixelMatrix [screenWidth][screenHeight]CellManipulator

func (spm *ScreenPixelMatrix) GetAllNonEmpty() []CellManipulator {
	var cells []CellManipulator

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

func (spm *ScreenPixelMatrix) ExportToLinkedList(ll *LinkedList) *LinkedList {
	// for _, col := range spm {
	// 	for _, cell := range col {
	// 		// cell should be pointer to CellManipulator value struct
	// 		if cell != nil {
	// 			ll.Add(cell.(NodeManipulator))
	// 		}
	// 	}
	// }
	for x := 0; x < len(spm); x++ {
		for y := 0; y < len(spm[x]); y++ {
			if spm[x][y] != nil {
				ll.Add(spm[x][y].(NodeManipulator))
			}
		}
	}

	return ll
}

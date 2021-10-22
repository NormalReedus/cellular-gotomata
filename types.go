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
	grid [screenWidth][screenHeight]CellManipulator
}

func NewGrid() *Grid {
	return &Grid{}
}

func (g *Grid) Bounds() (int, int) {
	return len(g.grid), len(g.grid[0])
}

func (g *Grid) Get(coords Point) CellManipulator {
	return g.grid[coords.X][coords.Y]
}

func (g *Grid) Move(currentCoords Point, newCoords Point, cell CellManipulator) {
	g.Remove(currentCoords)
	g.Set(newCoords, cell)
	cell.Position().SetCoords(newCoords.X, newCoords.Y)
}
func (g *Grid) Set(coords Point, cell CellManipulator) {
	g.grid[coords.X][coords.Y] = cell

	cell.SetParentGrid(g)
}

func (g *Grid) Remove(coords Point) {
	g.grid[coords.X][coords.Y] = nil
}

func (g *Grid) RandomOpenCell() (*Point, error) {
	var openCells []Point

	for x, column := range g.grid {
		for y := range column {
			if g.grid[x][y] == nil {
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

// Abstract struct that is embedded into Dot (i.e. not used directly anywhere)
// This makes any embedding struct implement the CellManipulator
type CellManipulator interface {
	SetParentGrid(grid *Grid)
	MoveCell(coords Point)
	RemoveCell()
	Position() *Point
}

type Cell struct {
	parentGrid *Grid
	position   Point
}

func (c *Cell) SetParentGrid(grid *Grid) {
	c.parentGrid = grid
}

func (c *Cell) MoveCell(coords Point) {
	c.parentGrid.Move(c.position, coords, c)
}

func (c *Cell) RemoveCell() {
	c.parentGrid.Remove(c.position)
}

func (c *Cell) Position() *Point {
	return &c.position
}

//* -------------------------
//* LINKED LIST
//* -------------------------

type LinkedList struct {
	head   NodeManipulator
	tail   NodeManipulator
	length int
}

// func NewLinkedList(nodes ...NodeManipulator) *LinkedList {
func NewLinkedList() *LinkedList {
	ll := &LinkedList{}

	// for _, node := range nodes {
	// 	ll.Add(node)
	// }

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

func (ll *LinkedList) SetHead(node NodeManipulator) NodeManipulator {
	ll.head = node
	return node
}

func (ll *LinkedList) Tail() NodeManipulator {
	return ll.tail
}

func (ll *LinkedList) SetTail(node NodeManipulator) NodeManipulator {
	ll.tail = node
	return node
}

func (ll *LinkedList) Length() int {
	return ll.length
}

func (ll *LinkedList) Add(node NodeManipulator) NodeManipulator {
	if ll.head == nil {
		ll.head = node
		ll.tail = node

		ll.incrementLength()

		return node
	}

	node.SetPrevNode(ll.tail)
	ll.tail.SetNextNode(node)
	ll.tail = node
	ll.incrementLength()

	node.SetParentList(ll)

	return node
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

		for {
			callback(node)

			// if node.prev == nil {
			if node.PrevNode() == nil {
				break
			}

			node = node.PrevNode()
		}

	} else {
		node := ll.head

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
	//TODO: this errors when right clicking to delete nodes
	// There are always 2 refs to delete to garbage collect this node...
	if n.prev == nil {
		//* If this node is head AND tail
		if n.next == nil {
			fmt.Println("NODE")
			fmt.Println(n)
			// Both refs are from list (head, tail), since there are no other nodes
			n.parentList.SetTail(nil)
			n.parentList.SetHead(nil)

			n.parentList.decrementLength()
			return
		}

		fmt.Println("NODE")
		fmt.Println(n)
		//* If this node is ONLY head
		// One ref from list (head) and one ref from next node (prev)
		n.parentList.SetHead(n.next)
		n.next.SetPrevNode(nil)

		n.parentList.decrementLength()
		return
	}

	//* If this node is ONLY tail
	if n.next == nil {
		fmt.Println("NODE")
		fmt.Println(n)
		// One ref from list (tail) and one ref from prev node (next)
		n.parentList.SetTail(n.prev)
		n.prev.SetNextNode(nil)

		n.parentList.decrementLength()
		return
	}

	//* If this node is NEITHER head nor tail
	// One ref from both prev (next) and next (prev)
	fmt.Println("NODE")
	fmt.Println(n)

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

func (p Point) String() string {
	return fmt.Sprintf("{ x: %d, y: %d }", p.X, p.Y)
}

func (p *Point) SetCoords(x, y int) {
	p.X = x
	p.Y = y
}
func (p *Point) GetCoords() (int, int) {
	return p.X, p.Y
}

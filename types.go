package main

import (
	"fmt"
	"image/color"
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

//* -------------------------
//* GRID
//* -------------------------

//TODO: continue with Grid and Cell the same as with LinkedList and Node, embedding Cell into Dot to implement the CellManipulator interface that together with Cell have the fields and methods to interact with the Grid
//TODO: grid should then use CellManipulator for everything, since Dot should implement CellManipulator. Grid should just be another way to reference Dots with O(1) lookups
//TODO: there should be something that adds a dot to both grid and ll at the same time (a method on Game?). It should take coords, create the Dot, call ll.Add and grid.Set with the coords
//TODO: maybe coords should be on Cell instead of Dot, see if there is a way to keep it the coords accessible through Position. This could be done by setting Position on Cell and embedding Cell on Dot

type Grid struct {
	grid [screenWidth][screenHeight]CellManipulator
}

func NewGrid() *Grid {
	return &Grid{}
}

func (dg *Grid) Get(x, y int) CellManipulator {
	return dg.grid[x][y]
}

func (dg *Grid) Set(x, y int, cell CellManipulator) {
	dg.grid[x][y] = cell
}

func (dg *Grid) Remove(x, y int) {
	dg.grid[x][y] = nil
}

// Abstract struct that is embedded into Dot (i.e. not used directly anywhere)
// This makes any embedding struct implement the CellManipulator
type CellManipulator interface {
	SetParentGrid(grid *Grid)
	// PrevNode() NodeManipulator
	// SetPrevNode(NodeManipulator) NodeManipulator
	// NextNode() NodeManipulator
	// SetNextNode(NodeManipulator) NodeManipulator
	RemoveCell()
	Position() *Point
}

type Cell struct {
	parentGrid *Grid
	position   *Point
}

func (c *Cell) SetParentGrid(grid *Grid) {
	c.parentGrid = grid
}

func (c *Cell) RemoveCell() {
	c.parentGrid.Remove(c.Position().X, c.Position().Y)
}

func (c *Cell) Position() *Point {
	return c.position
}

//* -------------------------
//* LINKED LIST
//* -------------------------

type LinkedList struct {
	head   NodeManipulator
	tail   NodeManipulator
	length int
}

func NewLinkedList(nodes ...NodeManipulator) *LinkedList {
	ll := &LinkedList{}

	for _, node := range nodes {
		ll.Add(node)
	}

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
	// There are always 2 refs to delete to garbage collect this node...
	if n.PrevNode() == nil {
		//* If this node is head AND tail
		if n.next == nil {
			// Both refs are from list (head, tail), since there are no other nodes
			n.parentList.SetTail(nil)
			n.parentList.SetHead(nil)

			n.parentList.decrementLength()
		}

		//* If this node is ONLY head
		// One ref from list (head) and one ref from next node (prev)
		n.parentList.SetHead(n.next)
		n.next.SetPrevNode(nil)

		n.parentList.decrementLength()
	}

	//* If this node is ONLY tail
	if n.next == nil {
		// One ref from list (tail) and one ref from prev node (next)
		n.parentList.SetTail(n.prev)
		n.prev.SetNextNode(nil)

		n.parentList.decrementLength()
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
	image        *ebiten.Image
	fill         color.Color
	screenBounds Point
	Node
	Cell
}

func NewDot(startX, startY int, screenWidth, screenHeight int) *Dot {
	image := ebiten.NewImage(1, 1)
	color, err := colorx.ParseHexColor("#adb5bd")
	if err != nil {
		log.Fatal(err)
	}

	return &Dot{
		image:        image,
		Cell:         Cell{position: &Point{startX, startY}},
		fill:         color,
		screenBounds: Point{screenWidth, screenHeight}, // TODO: Maybe not float64 for when we do comparisons later? I.e. don't use Point
	}
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
	return fmt.Sprintf("{ x: %d, y: %d }", int(p.X), int(p.Y))
}

func (p *Point) Set(x, y int) {
	p.X = x
	p.Y = y
}
func (p *Point) Get() (int, int) {
	return p.X, p.Y
}

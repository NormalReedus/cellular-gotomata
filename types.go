package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

//* -------------------------
//* GRID
//* -------------------------

type Grid struct {
	grid [screenWidth][screenHeight]*Dot
}

func NewGrid() *Grid {
	return &Grid{}
}

func (dg *Grid) Get(x, y int) *Dot {
	return dg.grid[x][y]
}

func (dg *Grid) Set(x, y int, dot *Dot) {
	dg.grid[x][y] = dot
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
		return fmt.Sprintf("LinkedList{ length: %d, head: %v, tail: %v }", ll.length, ll.head, ll.tail)
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
	SetParentList(list *LinkedList) *LinkedList
	PrevNode() NodeManipulator
	SetPrevNode(NodeManipulator) NodeManipulator
	NextNode() NodeManipulator
	SetNextNode(NodeManipulator) NodeManipulator
	RemoveNode() NodeManipulator
}

// Abstract struct that is embedded into Dot (i.e. not used directly anywhere)
// This makes any embedding struct implement the NodeManipulator
type Node struct {
	parentList *LinkedList
	prev       NodeManipulator
	next       NodeManipulator
}

// func NewNode(list *LinkedList) *Node {
// 	return &Node{list: list} // next and prev are set by LinkedList
// }

// func (n Node) String() string {
// 	var prev, next string

// 	if n.prev != nil {
// 		prev = n.prev.String()
// 	}
// 	if n.next != nil {
// 		next = n.next.String()
// 	}

// 	if prev == "" {
// 		prev = "<nil>"
// 	}
// 	if next == "" {
// 		next = "<nil>"
// 	}

// 	return fmt.Sprintf("DotNode{ data: %v, prev: %v, next: %v }", n.String(), prev, next)
// }

func (n *Node) SetParentList(list *LinkedList) *LinkedList {
	n.parentList = list
	return list
}

func (n *Node) PrevNode() NodeManipulator {
	return n.prev
}

func (n *Node) SetPrevNode(node NodeManipulator) NodeManipulator {
	n.prev = node
	return node
}

func (n *Node) NextNode() NodeManipulator {
	return n.next
}

func (n *Node) SetNextNode(node NodeManipulator) NodeManipulator {
	n.next = node
	return node
}

func (n *Node) RemoveNode() NodeManipulator {
	// There are always 2 refs to delete to garbage collect this node...
	if n.PrevNode() == nil {
		//* If this node is head AND tail
		if n.next == nil {
			// Both refs are from list (head, tail), since there are no other nodes
			n.parentList.SetTail(nil)
			n.parentList.SetHead(nil)

			n.parentList.decrementLength()
			return n
		}

		//* If this node is ONLY head
		// One ref from list (head) and one ref from next node (prev)
		n.parentList.SetHead(n.next)
		n.next.SetPrevNode(nil)

		n.parentList.decrementLength()
		return n
	}

	//* If this node is ONLY tail
	if n.next == nil {
		// One ref from list (tail) and one ref from prev node (next)
		n.parentList.SetTail(n.prev)
		n.prev.SetNextNode(nil)

		n.parentList.decrementLength()
		return n
	}

	//* If this node is NEITHER head nor tail
	// One ref from both prev (next) and next (prev)
	n.prev.SetNextNode(n.next)
	n.next.SetPrevNode(n.prev)

	n.parentList.decrementLength()
	return n
}

//* -------------------------
//* DOT
//* -------------------------

type Dot struct {
	image        *ebiten.Image
	fill         color.Color
	Position     Point
	screenBounds Point
	Node
}

func NewDot(startX, startY float64, screenWidth, screenHeight int) *Dot {
	image := ebiten.NewImage(1, 1)
	color, err := colorx.ParseHexColor("#adb5bd")
	if err != nil {
		log.Fatal(err)
	}

	return &Dot{
		image:        image,
		Position:     Point{startX, startY},
		fill:         color,
		screenBounds: Point{float64(screenWidth), float64(screenHeight)}, // TODO: Maybe not float64 for when we do comparisons later? I.e. don't use Point
	}
}

// Pretty print the Dot position x & y coordinates
func (d Dot) String() string {
	return fmt.Sprintf("Dot{ x: %d, y: %d, prev: %v, next: %v }", int(d.Position.X), int(d.Position.Y), d.prev, d.next)
}

func (d *Dot) Draw(screen *ebiten.Image) {
	dotOpts := &ebiten.DrawImageOptions{}
	dotOpts.GeoM.Translate(d.Position.X, d.Position.Y) // position
	d.image.Fill(d.fill)                               // color
	screen.DrawImage(d.image, dotOpts)
}

//* -------------------------
//* POINT
//* -------------------------
// Abstract struct that is embedded into Dot (i.e. not used directly anywhere)
type Point struct {
	X, Y float64
}

func (p *Point) Set(x, y float64) {
	p.X = x
	p.Y = y
}
func (p *Point) Get() (float64, float64) {
	return p.X, p.Y
}

package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/icza/gox/imagex/colorx"
)

//TODO: Don't embed Dot in DotNode, rather embed DotNode and DotCell inside Dot
// DotNode should have a parentList etc and all the other methods / fields for interacting with the LinkedDotList
// DotNode should have a parentGrid etc and all the other methods / field for interacting with the DotGrid
// Then DotGrid and LinkedDotList should use interfaces instead of *DotNode or *Dot etc
// Maybe the interfaces can be so generic, that we can rename the grid and ll to just Grid and LinkedList

type DotGrid struct {
	grid [screenWidth][screenHeight]*Dot
}

func NewDotGrid() *DotGrid {
	return &DotGrid{}
}

func (dg *DotGrid) Get(x, y int) *Dot {
	return dg.grid[x][y]
}

func (dg *DotGrid) Set(x, y int, dot *Dot) {
	dg.grid[x][y] = dot
}

type LinkedDotList struct {
	head   *DotNode
	tail   *DotNode
	length int
}

func NewLinkedDotList(dots ...*Dot) *LinkedDotList {
	ll := &LinkedDotList{}

	for _, dot := range dots {
		ll.Add(dot)
	}

	return ll
}

func (ll LinkedDotList) String() string {
	if ll.head != nil && ll.tail != nil {
		return fmt.Sprintf("LinkedDotList{ length: %d, head: %v, tail: %v }", ll.length, ll.head.Dot.String(), ll.tail.Dot.String())
	}

	return "LinkedDotList{ empty }"
}

func (ll *LinkedDotList) Head() *DotNode {
	return ll.head
}
func (ll *LinkedDotList) Tail() *DotNode {
	return ll.tail
}
func (ll *LinkedDotList) Length() int {
	return ll.length
}

func (ll *LinkedDotList) Add(dot *Dot) *DotNode {
	node := NewDotNode(dot, ll)

	if ll.head == nil {
		ll.head = node
		ll.tail = node

		ll.incrementLength()

		return node
	}

	node.prev = ll.tail
	ll.tail.next = node
	ll.tail = node
	ll.incrementLength()

	return node
}

func (ll *LinkedDotList) incrementLength() {
	ll.length++
}

func (ll *LinkedDotList) decrementLength() {
	ll.length--
}

// Loops through DotNodes, not Dots
// To access Dots you must use DotNode.data
func (ll LinkedDotList) ForEach(callback func(*DotNode), reverse bool) {
	if reverse {
		node := ll.tail

		for {
			callback(node)

			if node.prev == nil {
				break
			}

			node = node.prev
		}

	} else {
		node := ll.head

		for {
			callback(node)

			if node.next == nil {
				break
			}

			node = node.next
		}

	}
}

type DotNode struct {
	list *LinkedDotList
	*Dot // .Dot seemed excessively dotty to call all the time
	next *DotNode
	prev *DotNode
}

func NewDotNode(dot *Dot, list *LinkedDotList) *DotNode {
	return &DotNode{Dot: dot, list: list} // next and prev are set in LinkedDotList
}

func (dn DotNode) String() string {
	var prev, next string

	if dn.prev != nil {
		prev = dn.prev.Dot.String()
	}
	if dn.next != nil {
		next = dn.next.Dot.String()
	}

	if prev == "" {
		prev = "<nil>"
	}
	if next == "" {
		next = "<nil>"
	}

	return fmt.Sprintf("DotNode{ data: %v, prev: %v, next: %v }", dn.Dot.String(), prev, next)
}

func (dn *DotNode) Prev() *DotNode {
	return dn.prev
}

func (dn *DotNode) Next() *DotNode {
	return dn.next
}

func (dn *DotNode) Remove() *DotNode {
	// There are always 2 refs to delete to garbage collect this node...
	if dn.prev == nil {
		//* If this node is head AND tail
		if dn.next == nil {
			// Both refs are from list (head, tail), since there are no other nodes
			dn.list.tail = nil
			dn.list.head = nil

			dn.list.decrementLength()
			return dn
		}

		//* If this node is ONLY head
		// One ref from list (head) and one ref from next node (prev)
		dn.list.head = dn.next
		dn.next.prev = nil

		dn.list.decrementLength()
		return dn
	}

	//* If this node is ONLY tail
	if dn.next == nil {
		// One ref from list (tail) and one ref from prev node (next)
		dn.list.tail = dn.prev
		dn.prev.next = nil

		dn.list.decrementLength()
		return dn
	}

	//* If this node is NEITHER head nor tail
	// One ref from both prev (next) and next (prev)
	dn.prev.next = dn.next
	dn.next.prev = dn.prev

	dn.list.decrementLength()
	return dn
}

// Use NewDot constructor instead
type Dot struct {
	image        *ebiten.Image
	fill         color.Color
	Position     Point
	screenBounds Point
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
	return fmt.Sprintf("Dot{ x: %d, y: %d }", int(d.Position.X), int(d.Position.Y))
}

func (d *Dot) Draw(screen *ebiten.Image) {
	dotOpts := &ebiten.DrawImageOptions{}
	dotOpts.GeoM.Translate(d.Position.X, d.Position.Y) // position
	d.image.Fill(d.fill)                               // color
	screen.DrawImage(d.image, dotOpts)
}

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

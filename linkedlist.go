package main

//* -------------------------
//* LINKED LIST
//* -------------------------

// type LinkedList struct {
// 	head   NodeManipulator
// 	tail   NodeManipulator
// 	length int
// }

// func NewLinkedList() *LinkedList {
// 	ll := &LinkedList{}

// 	return ll
// }

// func NewLinkedListFromMatrix(matrix *ScreenPixelMatrix) *LinkedList {
// 	ll := &LinkedList{}

// 	ll = matrix.ExportToLinkedList(ll)

// 	return ll
// }

// func (ll LinkedList) String() string {
// 	if ll.head != nil && ll.tail != nil {
// 		return fmt.Sprintf("LinkedList{ nodeType: %s, length: %d, head: %v, tail: %v }", reflect.TypeOf(ll.head), ll.length, ll.head, ll.tail)
// 	}

// 	return "LinkedList{ empty }"
// }

// func (ll *LinkedList) Head() NodeManipulator {
// 	return ll.head
// }

// func (ll *LinkedList) SetHead(node NodeManipulator) {
// 	ll.head = node
// }

// func (ll *LinkedList) Tail() NodeManipulator {
// 	return ll.tail
// }

// func (ll *LinkedList) SetTail(node NodeManipulator) {
// 	ll.tail = node

// }

// func (ll *LinkedList) Length() int {
// 	return ll.length
// }

// func (ll *LinkedList) Add(node NodeManipulator) {
// 	if ll.head == nil {
// 		ll.head = node
// 		ll.tail = node
// 		ll.incrementLength()

// 		node.SetParentList(ll)

// 		return
// 	}

// 	node.SetPrevNode(ll.tail)

// 	ll.tail.SetNextNode(node)
// 	ll.tail = node
// 	ll.incrementLength()

// 	node.SetParentList(ll)
// }

// func (ll *LinkedList) incrementLength() {
// 	ll.length++
// }

// func (ll *LinkedList) decrementLength() {
// 	ll.length--
// }

// func (ll LinkedList) ForEach(callback func(node NodeManipulator), reverse bool) {
// 	if reverse {
// 		node := ll.tail

// 		// Don't do anything if there are no nodes
// 		if node == nil {
// 			return
// 		}

// 		for {
// 			callback(node)

// 			if node.PrevNode() == nil {
// 				break
// 			}

// 			node = node.PrevNode()
// 		}

// 	} else {
// 		node := ll.head

// 		// Don't do anything if there are no nodes
// 		if node == nil {
// 			return
// 		}

// 		for {
// 			callback(node)

// 			if node.NextNode() == nil {
// 				break
// 			}

// 			node = node.NextNode()
// 		}

// 	}
// }

//* -------------------------
//* NODE
//* -------------------------

// type NodeManipulator interface {
// 	SetParentList(list *LinkedList)
// 	PrevNode() NodeManipulator
// 	SetPrevNode(NodeManipulator)
// 	NextNode() NodeManipulator
// 	SetNextNode(NodeManipulator)
// 	RemoveNode()
// }

// // Abstract struct that is embedded into Dot (i.e. not used directly anywhere)
// // This makes any embedding struct implement the NodeManipulator
// type Node struct {
// 	parentList *LinkedList
// 	prev       NodeManipulator
// 	next       NodeManipulator
// }

// func (n *Node) String() string {
// 	// return fmt.Sprintf("Node{ parentList: %v, prev: %v, next: %v }", n.parentList, n.prev, n.next)
// 	return fmt.Sprintf("Node{ prev: %v, next: %v }", n.prev, n.next)
// }

// func (n *Node) SetParentList(list *LinkedList) {
// 	n.parentList = list
// }

// func (n *Node) PrevNode() NodeManipulator {
// 	return n.prev
// }

// func (n *Node) SetPrevNode(node NodeManipulator) {
// 	n.prev = node
// }

// func (n *Node) NextNode() NodeManipulator {
// 	return n.next
// }

// func (n *Node) SetNextNode(node NodeManipulator) {
// 	n.next = node
// }

// func (n *Node) RemoveNode() {
// 	// There are always 2 refs to delete to garbage collect this node...
// 	if n.prev == nil {
// 		//* If this node is head AND tail
// 		if n.next == nil {
// 			// Both refs are from list (head, tail), since there are no other nodes
// 			n.parentList.SetTail(nil)
// 			n.parentList.SetHead(nil)

// 			n.parentList.decrementLength()
// 			return
// 		}

// 		//* If this node is ONLY head
// 		// One ref from list (head) and one ref from next node (prev)
// 		n.parentList.SetHead(n.next)
// 		n.next.SetPrevNode(nil)

// 		n.parentList.decrementLength()
// 		return
// 	}

// 	//* If this node is ONLY tail
// 	if n.next == nil {
// 		// One ref from list (tail) and one ref from prev node (next)
// 		n.parentList.SetTail(n.prev)
// 		n.prev.SetNextNode(nil)

// 		n.parentList.decrementLength()
// 		return
// 	}

// 	//* If this node is NEITHER head nor tail
// 	// One ref from both prev (next) and next (prev)
// 	n.prev.SetNextNode(n.next)
// 	n.next.SetPrevNode(n.prev)

// 	n.parentList.decrementLength()
// }

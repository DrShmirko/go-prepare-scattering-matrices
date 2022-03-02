package doublyLinkedList

import (
	"github.com/kshmirko/prepare-mueller-matrices/calcresult"
)

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{}
}

func (d *DoublyLinkedList) PushFront(data *calcresult.CalculusResult) {

	newnode := &Node{data: data}

	if d.head == nil {
		d.head = newnode
		d.tail = newnode
	} else {
		newnode.next = d.head
		d.head.prev = newnode
		d.head = newnode
	}
	d.len++
	return
}

func (d *DoublyLinkedList) PushBack(data *calcresult.CalculusResult) {
	newnode := &Node{
		data: data,
	}

	if d.head == nil {
		d.head = newnode
		d.tail = newnode
	} else {
		currentNode := d.head
		for currentNode.next != nil {
			currentNode = currentNode.next
		}
		newnode.prev = currentNode
		currentNode.next = newnode
		d.tail = newnode
	}
	d.len++

	return
}

func (d *DoublyLinkedList) Front() *Node {
	return d.head
}

func (d *DoublyLinkedList) Tail() *Node {
	return d.tail
}

func (d *DoublyLinkedList) Size() int {
	return d.len
}

func (n *Node) Value() *calcresult.CalculusResult {
	return n.data
}

func (n *Node) Next() *Node {
	return n.next
}

func (n *Node) Prev() *Node {
	return n.prev
}

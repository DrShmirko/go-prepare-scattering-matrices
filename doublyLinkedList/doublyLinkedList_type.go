package doublyLinkedList

import (
	"github.com/kshmirko/prepare-mueller-matrices/calcresult"
)

type Node struct {
	data       *calcresult.CalculusResult
	prev, next *Node
}

type DoublyLinkedList struct {
	len  int
	tail *Node
	head *Node
}

package doublyLinkedList

import (
	"github.com/kshmirko/prepare-mueller-matrices/calcresult"
)

// Node - описание узла списка. Содержит указатели на следующий и предыдущий
// элементы
type Node struct {
	data       *calcresult.CalculusResult
	prev, next *Node
}

// DoublyLinkedList - структура списка, управляет самими спискаом. содержит указатели
// на хвост и голову списка, его длину
type DoublyLinkedList struct {
	len  int
	tail *Node
	head *Node
}

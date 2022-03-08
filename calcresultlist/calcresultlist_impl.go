package calcresultlist

import (
	"fmt"

	"github.com/kshmirko/prepare-mueller-matrices/calcresult"
	"github.com/kshmirko/prepare-mueller-matrices/doublyLinkedList"
)

// NewCalcResultsList - конструктор нового специализированного типа для
// хранения результатов расчетов
func NewCalcResultsList(prefix string) *CalcResultsList {
	return &CalcResultsList{
		DoublyLinkedList: doublyLinkedList.NewDoublyLinkedList(),
		Prefix:           prefix,
	}
}

// ApplyForward - apply function fun to every element of a list from the first element to the end.
func (c *CalcResultsList) ApplyForward(fun func(cr *calcresult.CalculusResult)) error {

	if c.Front() == nil {
		return fmt.Errorf("applyforward error: list is empty")
	}
	tmp := c.Front()
	for tmp != nil {
		p := tmp.Value()
		fun(p)
		tmp = tmp.Next()
	}
	return nil
}

// ApplyBackward - apply function fun to every element of a list from the end element to the first.
func (c *CalcResultsList) ApplyBackward(fun func(cr *calcresult.CalculusResult)) error {

	if c.Front() == nil {
		return fmt.Errorf("applybackward error: list is empty")
	}
	tmp := c.Tail()
	for tmp != nil {
		p := tmp.Value() //.(*CalculusResult)
		fun(p)
		tmp = tmp.Prev()
	}
	return nil
}

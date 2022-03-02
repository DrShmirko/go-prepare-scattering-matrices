package calcresultlist

import (
	"github.com/kshmirko/prepare-mueller-matrices/doublyLinkedList"
)

type CalcResultsList struct {
	*doublyLinkedList.DoublyLinkedList
	prefix string
}

// Глобальные переменные хранят  соответствующие списки результатов расчета
var (
	SpheroidList *CalcResultsList
	SpheresList  *CalcResultsList
	CombList     *CalcResultsList
)

func init() {
	SpheroidList = NewCalcResultsList("sphrd")
	SpheresList = NewCalcResultsList("sphrs")
	CombList = NewCalcResultsList("total")
}
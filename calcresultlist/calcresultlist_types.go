package calcresultlist

import (
	"github.com/kshmirko/prepare-mueller-matrices/doublyLinkedList"
)

// CalcResultsList - список для хранения результатов расчета. Включает в себя
// двусвязный список и строковую пееменную для идентификации данных
type CalcResultsList struct {
	*doublyLinkedList.DoublyLinkedList
	prefix string
}

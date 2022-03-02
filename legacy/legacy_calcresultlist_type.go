package legacy

import "container/list"

// CalcResultsList Для чего используем список
type CalcResultsList struct {
	*list.List
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

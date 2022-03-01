package legacy

import "container/list"

// CalcResultsList Для чего используем список
type CalcResultsList struct {
	*list.List
	prefix string
}

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

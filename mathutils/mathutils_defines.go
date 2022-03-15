package mathutils

// // PolyFit - структура для хранения полиномиальной регрессии
// type PolyFit struct {
// 	order  int
// 	x, y   []float64
// 	Coeffs []float64
// }

// // IFit - интерфейс для различного рода математических фильтров, сейчас от отвечает
// // за полиномиальный подгон
// type IFit interface {
// 	Fit() error
// 	SetXY(ax, ay []float64)
// 	Evaluate(x float64) (ret float64)
// 	EvaluateArray(x []float64) (ret []float64)
// }

// Function is a piecewise-linear 1-dimensional function
type Function struct {
	X []float64
	Y []float64
}

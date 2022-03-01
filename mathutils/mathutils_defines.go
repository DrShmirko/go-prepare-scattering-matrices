package mathutils

type PolyFit struct {
	order  int
	x, y   []float64
	Coeffs []float64
}

type IPolyFit interface {
	Fit() error
	SetXY(ax, ay []float64)
	Evaluate(x float64) (ret float64)
	EvaluateArray(x []float64) (ret []float64)
}

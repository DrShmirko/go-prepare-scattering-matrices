package mathutils

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func NewPolyFit(order int) IFit {
	return &PolyFit{
		order:  order,
		x:      nil,
		y:      nil,
		Coeffs: nil,
	}
}

// Fit - решаем систему уравнений, ищем вектор решений - он и будет нашими коффициентами
// подгона
func (p *PolyFit) Fit() error {
	p.Coeffs = make([]float64, 1+p.order)

	a := vandermonde(p.x, p.order+1)
	b := mat.NewDense(len(p.y), 1, p.y)
	c := mat.NewDense(p.order+1, 1, p.Coeffs)

	var qr mat.QR
	qr.Factorize(a)
	const trans = false
	err := qr.SolveTo(c, trans, b)
	if err != nil {
		fmt.Printf("could not solve QR: %+v", err)
	}
	return err
}

// SetXY -  функция для замены исходных данных без повторного создания объекта
func (p *PolyFit) SetXY(ax, ay []float64) {
	p.x = append([]float64(nil), ax...)
	p.y = append([]float64(nil), ay...)
}

// Evaluate - вычисляет значение  полинома по заданным коэффициентам в выбранной точке
func (p *PolyFit) Evaluate(x float64) (ret float64) {
	ret = 0
	tmp := 1.0
	for _, v := range p.Coeffs {
		ret = ret + v*tmp
		tmp = tmp * x
	}
	return
}

// EvaluateArray - вычисляет значения аппроксимирующего полинома для каждого из
// значений массива
func (p *PolyFit) EvaluateArray(x []float64) (ret []float64) {
	ret = make([]float64, len(x))
	for i, xi := range x {
		ret[i] = p.Evaluate(xi)
	}
	return
}

// vandermonde - матрица Вандермонда для вычисления коэффициентов полиномиального
// подгона
// Вид матрицы:
// ============
// | x0^0, x0^1, x0^2, ..., x0^m-1 |
// | x1^0, x1^1, x1^2, ..., x1^m-1 |
// | x2^0, x2^1, x2^2, ..., x2^m-1 |
// | x3^0, x3^1, x3^2, ..., x3^m-1 |
// | ............................. |
// | xn^0, xn^1, xn^2, ..., xn^m-1 |
// ---------------------------------
func vandermonde(a []float64, d int) *mat.Dense {
	x := mat.NewDense(len(a), d, nil)
	for i := range a {
		for j, p := 0, 1.0; j < d; j, p = j+1, p*a[i] {
			x.Set(i, j, p)
		}
	}
	return x
}

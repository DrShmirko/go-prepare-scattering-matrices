package mathutils

import "math"

// func NewPolyFit(order int) IFit {
// 	return &PolyFit{
// 		order:  order,
// 		x:      nil,
// 		y:      nil,
// 		Coeffs: nil,
// 	}
// }

// Fit - решаем систему уравнений, ищем вектор решений - он и будет нашими коффициентами
// подгона
// func (p *PolyFit) Fit() error {
// 	p.Coeffs = make([]float64, 1+p.order)

// 	a := vandermonde(p.x, p.order+1)
// 	b := mat.NewDense(len(p.y), 1, p.y)
// 	c := mat.NewDense(p.order+1, 1, p.Coeffs)

// 	var qr mat.QR
// 	qr.Factorize(a)
// 	const trans = false
// 	err := qr.SolveTo(c, trans, b)
// 	if err != nil {
// 		fmt.Printf("could not solve QR: %+v", err)
// 	}
// 	return err
// }

// // SetXY -  функция для замены исходных данных без повторного создания объекта
// func (p *PolyFit) SetXY(ax, ay []float64) {
// 	p.x = append([]float64(nil), ax...)
// 	p.y = append([]float64(nil), ay...)
// }

// // Evaluate - вычисляет значение  полинома по заданным коэффициентам в выбранной точке
// func (p *PolyFit) Evaluate(x float64) (ret float64) {
// 	ret = 0
// 	tmp := 1.0
// 	for _, v := range p.Coeffs {
// 		ret = ret + v*tmp
// 		tmp = tmp * x
// 	}
// 	return
// }

// // EvaluateArray - вычисляет значения аппроксимирующего полинома для каждого из
// // значений массива
// func (p *PolyFit) EvaluateArray(x []float64) (ret []float64) {
// 	ret = make([]float64, len(x))
// 	for i, xi := range x {
// 		ret[i] = p.Evaluate(xi)
// 	}
// 	return
// }

// // vandermonde - матрица Вандермонда для вычисления коэффициентов полиномиального
// // подгона
// // Вид матрицы:
// // ============
// // | x0^0, x0^1, x0^2, ..., x0^m-1 |
// // | x1^0, x1^1, x1^2, ..., x1^m-1 |
// // | x2^0, x2^1, x2^2, ..., x2^m-1 |
// // | x3^0, x3^1, x3^2, ..., x3^m-1 |
// // | ............................. |
// // | xn^0, xn^1, xn^2, ..., xn^m-1 |
// // ---------------------------------
// func vandermonde(a []float64, d int) *mat.Dense {
// 	x := mat.NewDense(len(a), d, nil)
// 	for i := range a {
// 		for j, p := 0, 1.0; j < d; j, p = j+1, p*a[i] {
// 			x.Set(i, j, p)
// 		}
// 	}
// 	return x
// }

// Area returns the definite integral of the function on its domain X.
//
// Time complexity: O(N), where N is the number of points.
// Space complexity: O(1)
func (f Function) Area() (area float64) {
	X, Y := f.X, f.Y
	for i := 1; i < len(X); i++ {
		area += (X[i] - X[i-1]) * (Y[i] + Y[i-1]) / 2
	}
	return area
}

// AreaUpTo returns the definite integral of the function on its domain X intersected with [-Inf, x].
//
// Time complexity: O(N), where N is the number of points.
// Space complexity: O(1)
func (f Function) AreaUpTo(x float64) (area float64) {
	X, Y := f.X, f.Y
	for i := 1; i < len(X); i++ {
		dX := X[i] - X[i-1]
		if x < X[i] {
			if x >= X[i-1] {
				dxX := x - X[i-1]
				w := dxX / dX
				y := (1-w)*Y[i-1] + w*Y[i]
				area += dxX * (y + Y[i-1]) / 2
			}
			return area
		}
		area += dX * (Y[i] + Y[i-1]) / 2
	}
	return area
}

// IsInterpolatedAt returns true if x is within the given range of points, false if outside of that range
func (f Function) IsInterpolatedAt(x float64) bool {
	n := len(f.X)
	if n == 0 {
		return false
	}
	left, right := f.X[0], f.X[n-1]
	return x >= left && x <= right
}

// At returns the function's value at the given point.
// Outside its domain X, the function is constant at 0.
//
// The function's X and Y slices are expected to be the same legnth. The length property is _not_ verified.
// The function's X slice is expected to be sorted in ascending order. The sortedness property is _not_ verified.
//
// Time complexity: O(log(N)), where N is the number of points.
// Space complexity: O(1)
func (f Function) At(x float64) float64 {
	X, Y := f.X, f.Y
	i, j := 0, len(X)
	for i < j {
		h := int(uint(i+j) >> 1)
		if X[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	if i == 0 {
		if len(X) > 0 && x < X[0] {
			return 0
		}
		if len(X) > 0 && x == X[0] {
			return Y[0]
		}
		return 0
	}
	if i == len(X) {
		return 0
	}
	w := (x - X[i-1]) / (X[i] - X[i-1])
	return (1-w)*Y[i-1] + w*Y[i]
}

// Span generates `nPoints` equidistant points spanning [min,max]
func Span(min, max float64, nPoints int) []float64 {
	X := make([]float64, nPoints)
	min, max = math.Min(max, min), math.Max(max, min)
	d := max - min
	for i := range X {
		X[i] = min + d*(float64(i)/float64(nPoints-1))
	}
	return X
}

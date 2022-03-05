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

func (p *PolyFit) SetXY(ax, ay []float64) {
	p.x = append([]float64(nil), ax...)
	p.y = append([]float64(nil), ay...)
}

func (p *PolyFit) Evaluate(x float64) (ret float64) {
	ret = 0
	tmp := 1.0
	for _, v := range p.Coeffs {
		ret = ret + v*tmp
		tmp = tmp * x
	}
	return
}

func (p *PolyFit) EvaluateArray(x []float64) (ret []float64) {
	ret = make([]float64, len(x))
	for i, xi := range x {
		ret[i] = p.Evaluate(xi)
	}
	return
}

func vandermonde(a []float64, d int) *mat.Dense {
	x := mat.NewDense(len(a), d, nil)
	for i := range a {
		for j, p := 0, 1.0; j < d; j, p = j+1, p*a[i] {
			x.Set(i, j, p)
		}
	}
	return x
}

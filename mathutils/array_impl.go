package mathutils

import (
	"fmt"
	"log"
)

func NewSimpleMatrix(r, c int, d []float64) *SimpleMatrix {
	if r <= 0 || c <= 0 {
		log.Fatalf("Rows and Cols must be >0")
	}

	var tmpd []float64
	if d == nil {
		tmpd = make([]float64, r*c)
	} else {
		tmpd = d
	}

	if r*c < len(d) {
		log.Fatalf("Length of d shoukd be greater or equat to Rows*Cols")
	}

	return &SimpleMatrix{
		Rows: r,
		Cols: c,
		Data: tmpd[:r*c],
	}
}

func (sm *SimpleMatrix) At(i, j int) float64 {
	if i < 0 || j < 0 {
		log.Fatalf("At: i and j must be >=0, %d, %d", i, j)
	}
	if i >= sm.Rows || j >= sm.Cols {
		log.Fatalf("i and j must less than Rosw and Cols respectively")
	}

	return sm.Data[i*sm.Cols+j]
}

func (sm *SimpleMatrix) Set(i, j int, v float64) {
	if i < 0 || j < 0 {
		log.Fatalf("Set: i and j must be >=0, %d, %d", i, j)
	}
	if i >= sm.Rows || j >= sm.Cols {
		log.Fatalf("i and j must less than Rosw and Cols respectively")
	}

	sm.Data[i*sm.Cols+j] = v
}

func (sm *SimpleMatrix) SetCol(j int, v []float64) error {
	if j >= sm.Cols || j < 0 {
		return fmt.Errorf("j out of bounds")
	}

	if sm.Rows != len(v) {
		return fmt.Errorf("num of elemendt in v and number of rows in sm should be equal")
	}

	for i, t := range v {
		sm.Data[i*sm.Cols+j] = t
	}

	return nil
}

func (sm *SimpleMatrix) Dims() (int, int) {
	return sm.Rows, sm.Cols
}

func (sm *SimpleMatrix) Scale(a float64) *SimpleMatrix {

	ret := &SimpleMatrix{
		Rows: sm.Rows,
		Cols: sm.Cols,
		Data: make([]float64, sm.Rows*sm.Cols),
	}

	for i, v := range sm.Data {
		ret.Data[i] = a * v
	}
	return ret
}

func (sm *SimpleMatrix) Add(b *SimpleMatrix) *SimpleMatrix {
	if (sm.Rows != b.Rows) || (sm.Cols != b.Cols) {
		log.Fatal("matrix dimension dhould be the same")
	}

	for i := range sm.Data {
		sm.Data[i] = sm.Data[i] + b.Data[i]
	}
	return sm
}

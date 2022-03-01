package legacy

import "gonum.org/v1/gonum/mat"

// CalculusResult - результаты расчетов
type CalculusResult struct {
	isSpheroid                    bool
	RecordId                      int
	Ext, Sca, Absb, VolC, Lr, MuL float64
	MuellerMat                    *mat.Dense
	Angle                         []float64
	SphericalFraction             float64
}

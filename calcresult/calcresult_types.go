package calcresult

import (
	"gonum.org/v1/gonum/mat"
	"time"
)

// CalculusResult - результаты расчетов, описание того, что получено
// программой
type CalculusResult struct {
	isSpheroid                    bool
	RecordId                      int
	Ext, Sca, Absb, VolC, Lr, MuL float64
	MuellerMat                    *mat.Dense
	Angle                         []float64
	SphericalFraction             float64
	Dt                            time.Time
}

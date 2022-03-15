package calcresult

import (
	"time"

	"github.com/kshmirko/prepare-mueller-matrices/mathutils"
)

// CalculusResult - результаты расчетов, описание того, что получено
// программой
type CalculusResult struct {
	isSpheroid                    bool
	RecordId                      int
	Ext, Sca, Absb, VolC, Lr, MuL float64
	MuellerMat                    *mathutils.SimpleMatrix
	Angle                         []float64
	SphericalFraction             float64
	Dt                            time.Time
}

// agglomerated - пакет для подготовки данных расчетов
// Загружает матрицы Мюллера из фала netcdf делает интрерполяцию к искомым
// параметрам mre, mim, xsize
package agglomerated

import (
	"github.com/kshmirko/prepare-mueller-matrices/mathutils"
)

type SingleParticle struct {
	ReIdx, ImIdx, Xsize, Angle []float64
	SM                         []*mathutils.SimpleMatrix
}

package scattlib

import (
	"github.com/kshmirko/prepare-mueller-matrices/legacy"
)

// MuellerMatrixAERONET - описывает вснутреннюю структурв данных,
// необходимых для расчета матриц рассеяния по данным AERONET
type MuellerMatrixAERONET struct {
	Wvl       float64
	IsSpheres bool
	dll       *legacy.Singleton
}

// DEBUG Для отладки. True - добавляет вывод отладочной информаци
const DEBUG = false

// aeronet - Ок, мы разместим здесь описание структуры файла аэронет
package aeronet

import "time"

type AeronetDataset struct {
	TimePoint         time.Time
	SphericalFraction float64
	ReIdx, ImIdx      []float64
	VolSd             []float64
	AotExt            []float64
}

type AeronetDatasets []AeronetDataset

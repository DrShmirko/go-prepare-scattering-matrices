package calcresult

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func NewCalcResult(isspheroid bool, rec_id int, ext, sca, absb, volc, lr, mul float64,
	mm *mat.Dense, angle []float64, sf float64) *CalculusResult {
	return &CalculusResult{
		isSpheroid:        isspheroid,
		RecordId:          rec_id,
		Ext:               ext,
		Sca:               sca,
		Absb:              absb,
		VolC:              volc,
		Lr:                lr,
		MuL:               mul,
		MuellerMat:        mm,
		Angle:             angle,
		SphericalFraction: sf,
	}
}

// DoPlotPolarization Рисует график угловой зависимости степени линейной поляризции
// и сохраняет его в файл saveto. Формат файла выбирается на основании расширения.
// Поддерживаются форматы png, pdf.
// Возвращает nil в случае успешного завершения метода
func (c *CalculusResult) DoPlotPolarization(saveto string) error {
	angle := c.Angle
	mat := c.MuellerMat
	pts := make(plotter.XYs, len(angle))

	// MuL = - S12/S11 * 100
	// Но S11 на самом деле S11
	for i, v := range angle {
		S11 := mat.At(i, 0)
		S12 := mat.At(i, 1)
		pts[i].X = v
		pts[i].Y = -S12 / S11 * 100.0
	}

	plt := plot.New()
	sf := fmt.Sprintf("SF = %4.2f", c.SphericalFraction)
	err := plotutil.AddLines(plt, sf, pts)
	plt.X.Label.Text = "Scattering angle"
	plt.Y.Label.Text = "Degree of linear polarization"
	plt.Title.Text = "μL vs ϑ"

	if err != nil {
		return err
	}

	plt.Add(plotter.NewGrid())
	plt.Legend.Top = true
	plt.X.Padding = 0 * vg.Centimeter
	plt.Y.Padding = 0 * vg.Centimeter
	err = plt.Save(15*vg.Centimeter, 12*vg.Centimeter, saveto)

	if err != nil {
		return err
	}

	return nil
}

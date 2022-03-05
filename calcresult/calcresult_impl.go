package calcresult

import (
	"github.com/kshmirko/prepare-mueller-matrices/plotdumper"
	"github.com/vdobler/chart"
	"gonum.org/v1/gonum/mat"
	"image/color"
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
// и сохраняет его в файл saveto. Сохоаняются в файлы png, svg и txt
// Возвращает nil в случае успешного завершения метода
//                                   Polarization
//
//
//       50  +
//           |                                 ###
//           |       :                        ## ###
//           |       :                      ##     ##
//       40  +       :                     ##       ##
//           |       :                    ##         #
//           |       :                   ##           #
//       30  +       :                  ##            ##
//           |       :                 ##              #
//           |       :                 #                #
//           |       :                #                 ##
//  D    20  +       :               ##                  #
//  O        |       :              #                    ##
//  P        |       :             #                      #
//       10  +       :           ##                       ##
//           |       :         ###                         #
//           |       :       ###                            #       #
//           |       :  ######                              #       ##
//        0  +- - - -#### - - - - - - - - - - - - - - - - - -#- - -#-## - - - -
//           |       :                                       #     # ##
//           |       :                                        #   #   #
//      -10  +       :                                         #  #
//           |       :                                         ####
//           |       :                                          #
//      -20  +-------+-------+-------+--------+-------+-------+-------+--------+
//          -30      0      30      60       90      120     150     180      210
//                                    Scattering angle
func (c *CalculusResult) DoPlotPolarization(saveto string) error {
	dumper := plotdumper.NewPlotDumper(saveto, 1, 1, 800, 600)
	defer dumper.Close()

	angle := c.Angle
	mat := c.MuellerMat
	dp := make([]float64, len(angle))

	for i, _ := range angle {
		S11 := mat.At(i, 0)
		S12 := mat.At(i, 1)
		dp[i] = -S12 / S11 * 100.0
	}

	pl := chart.ScatterChart{Title: "Polarization"}

	pl.XRange.Label, pl.YRange.Label = "Scattering angle", "DOP"
	pl.XRange.TicSetting.Grid = 1
	pl.AddDataPair("", angle, dp, chart.PlotStyleLines,
		chart.Style{Symbol: '#', SymbolColor: color.NRGBA{0x00, 0x00, 0xff, 0xff}, LineStyle: chart.SolidLine})
	pl.XRange.ShowZero = true
	pl.YRange.ShowZero = true
	pl.XRange.Fixed(0, 180, 30)
	pl.XRange.TicSetting.Delta = 30
	pl.YRange.TicSetting.Grid = 1

	dumper.Plot(&pl)
	return nil
}

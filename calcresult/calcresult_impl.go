package calcresult

import (
	"fmt"
	"image/color"
	"os"

	"github.com/kshmirko/prepare-mueller-matrices/mathutils"
	"github.com/kshmirko/prepare-mueller-matrices/plotdumper"
	"github.com/vdobler/chart"
)

func NewCalcResult(isspheroid bool, recId int, ext, sca, absb, volc, lr, mul float64,
	mm *mathutils.SimpleMatrix, angle []float64, sf float64) *CalculusResult {
	return &CalculusResult{
		isSpheroid:        isspheroid,
		RecordId:          recId,
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
	dumper, err := plotdumper.NewPlotDumper(saveto, 1, 1, 800, 600)

	if err != nil {
		return err
	}
	defer func() {
		err := dumper.Close()
		if err != nil {
			fmt.Print(err)
		}
	}()

	angle := c.Angle
	mtrx := c.MuellerMat
	dp := make([]float64, len(angle))

	for i := range angle {
		S11 := mtrx.At(i, 0)
		S12 := mtrx.At(i, 1)
		dp[i] = -S12 / S11 * 100.0
	}

	pl := chart.ScatterChart{Title: "Polarization"}

	pl.XRange.Label, pl.YRange.Label = "Scattering angle", "DOP, %"
	pl.XRange.TicSetting.Grid = 1
	pl.AddDataPair("", angle, dp, chart.PlotStyleLines,
		chart.Style{Symbol: '#', SymbolColor: color.NRGBA{B: 0xff, A: 0xff}, LineStyle: chart.SolidLine})
	pl.XRange.ShowZero = true
	pl.YRange.ShowZero = true
	pl.XRange.Fixed(0, 180, 30)
	pl.XRange.TicSetting.Delta = 30
	pl.YRange.TicSetting.Grid = 1

	err = dumper.Plot(&pl)
	if err != nil {
		return err
	}
	return nil
}

// DoSaveMatrixToFile - сохраняет матрицы в требуемом формате в файл fname
// в случае возниконовения и
func (c *CalculusResult) DoSaveMatrixToFile(fname string) error {

	fout, err := os.Create(fname)

	if err != nil {
		return fmt.Errorf("error creating file %s, %s", fname, err)
	}

	defer func() {
		err := fout.Close()
		if err != nil {
			fmt.Printf("Error closing file %s, err=%s", fname, err)
		}
	}()

	angle := c.Angle
	M := c.MuellerMat
	rows, _ := M.Dims()

	_, _ = fmt.Fprintf(fout, "%9.3e\t%9.3e\n", c.Sca/c.Ext, c.Ext)
	_, _ = fmt.Fprintf(fout, "%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t"+
		"%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\n",
		"Angle", "S11", "S12", "S13", "S14", "S21", "S22", "S23",
		"S24", "S31", "S32", "S33", "S34", "S41", "S42", "S43",
		"S44")
	_, _ = fmt.Fprintf(fout, "%d\n", len(angle))
	for i := 0; i < rows; i++ {
		_, _ = fmt.Fprintf(fout, "%9.3f\t", angle[i])
		_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 0))  //S11
		_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 1))  //S12
		_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)         //S13
		_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)         //S13
		_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 1))  //S21
		_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 2))  //S22
		_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)         //S23
		_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)         //S24
		_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)         //S31
		_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)         //S32
		_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 3))  //S33
		_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 4))  //S34
		_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)         //S41
		_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)         //S42
		_, _ = fmt.Fprintf(fout, "%9.3e\t", -M.At(i, 4)) //S43
		_, _ = fmt.Fprintf(fout, "%9.3e\n", M.At(i, 5))  //S44
	}
	return nil
}

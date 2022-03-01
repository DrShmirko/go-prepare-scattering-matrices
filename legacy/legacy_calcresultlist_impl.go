package legacy

import (
	"container/list"
	"fmt"
	"os"
)

// NewCalcResultsList - конструктор нового специализированного типа для хранения
// результатов расчетов
func NewCalcResultsList(prefix string) *CalcResultsList {
	return &CalcResultsList{
		List:   list.New(),
		prefix: prefix,
	}
}

func (c *CalcResultsList) Apply(fun func(cr *CalculusResult)) error {
	errCnt := 0
	for tmp := c.Front(); tmp != nil; tmp = tmp.Next() {
		p, ok := tmp.Value.(*CalculusResult)
		if ok {
			fun(p)
		} else {
			errCnt++
		}
	}
	if errCnt != 0 {
		return fmt.Errorf("при сохранении данных %d записи имеют странный тип", errCnt)
	}
	return nil
}

// SaveResults - сохрняем весь список, каждый элемент в свой собственный файл
func (c *CalcResultsList) SaveResults() error {

	err := c.Apply(func(cr *CalculusResult) {
		fname := fmt.Sprintf("out/%s_%05d.out", c.prefix, cr.RecordId)
		fout, err := os.Create(fname)

		if err != nil {
			fmt.Printf("Error creating file %s, %s\n", fname, err)
			return
		}

		defer func() {
			err := fout.Close()
			if err != nil {

			}
		}()

		angle := cr.Angle
		M := cr.MuellerMat
		rows, _ := M.Dims()
		Vc := cr.VolC
		_, _ = fmt.Fprintf(fout, "%9.3e\t%9.3e\t%9.3e\t%9.3e\t%9.3e\t%9.3e# Ext/V, Sca/V, Absb/V, LR, MuL, V \n",
			cr.Ext/Vc, cr.Sca/Vc, cr.Absb/Vc, cr.Lr, cr.MuL, cr.VolC)
		_, _ = fmt.Fprintf(fout, "%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t"+
			"%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\n",
			"Angle", "S11", "S12", "S13", "S14", "S21", "S22", "S23",
			"S24", "S31", "S32", "S33", "S34", "S41", "S42", "S43",
			"S44")
		for i := 0; i < rows; i++ {
			_, _ = fmt.Fprintf(fout, "%9.3f\t", angle[i])
			_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 0)/Vc)  //S11
			_, _ = fmt.Fprintf(fout, "%9.3e\t", -M.At(i, 1)/Vc) //S12
			_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)            //S13
			_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)            //S13
			_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 1)/Vc)  //S21
			_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 2)/Vc)  //S22
			_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)            //S23
			_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)            //S24
			_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)            //S31
			_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)            //S32
			_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 3)/Vc)  //S33
			_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 4)/Vc)  //S34
			_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)            //S41
			_, _ = fmt.Fprintf(fout, "%9.3e\t", 0.0)            //S42
			_, _ = fmt.Fprintf(fout, "%9.3e\t", -M.At(i, 4)/Vc) //S43
			_, _ = fmt.Fprintf(fout, "%9.3e\n", M.At(i, 5)/Vc)  //S44
		}

		return
	})

	return err
}

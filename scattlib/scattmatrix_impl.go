package scattlib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/kshmirko/prepare-mueller-matrices/aeronet"
	"github.com/kshmirko/prepare-mueller-matrices/calcresult"
	"github.com/kshmirko/prepare-mueller-matrices/calcresultlist"
	"github.com/kshmirko/prepare-mueller-matrices/legacy"
	"gonum.org/v1/gonum/interp"
	"gonum.org/v1/gonum/mat"
)

// NewMuellerMatrixAERONET конструктор объекта. Создаем экземпляр типа, загружаем библиотеку,
// загружаем настройки
func NewMuellerMatrixAERONET(wvl float64) *MuellerMatrixAERONET {
	instance := &MuellerMatrixAERONET{
		Wvl: wvl,
		dll: legacy.Instance(),
	}

	instance.dll.ReadConfig(true)
	instance.dll.SetWl(wvl)
	instance.dll.AllocateMemory()

	return instance
}

// Run - производит вычисления над заданным файлом
func (a *MuellerMatrixAERONET) Run(fname string, sf float64, skiprows int, matdir, picdir string) {
	// Read csv file into dataframe
	content, _ := ioutil.ReadFile(fname)
	ioContent := bufio.NewReader(strings.NewReader(string(content)))

	_ = sf
	// Skip necessary amount of  lines
	for i := 0; i < skiprows; i++ {
		_, err := ioContent.ReadString('\n')
		if err != nil {
			log.Printf("При чтении файла возникла ошибка (возможно файл пуст). %v", err)
		}
	}

	// Читаем файл в DataFrame. Важно, чтобы структура таблицы данных
	// сохраняла свой формат, я имею ввиду, сделование столбцов с данными
	df := dataframe.ReadCSV(ioContent, dataframe.HasHeader(true))

	// подготавливаем инструмент для аппроксимации оптической толщи
	//pf := mathutils.NewPolyFit(1)

	// Prepare for interpolation of refractive index
	pl := interp.PiecewiseLinear{}

	// Allocate memory for wavelength arrays
	LgWavelength := make([]float64, len(aeronet.IdxAot))
	LgAerOptThickness := make([]float64, len(aeronet.IdxAot))

	// Initialize lists for storage processed items
	SpheroidList := calcresultlist.NewCalcResultsList("sphrd")
	SpheresList := calcresultlist.NewCalcResultsList("sphrs")
	CombList := calcresultlist.NewCalcResultsList("total")

	// Before Rapply, let's filter out those records, in which Spherical Fraction < sf.
	df = df.Filter(dataframe.F{
		Colidx:     aeronet.IdxSphrericalFract,
		Comparator: series.LessEq,
		Comparando: sf,
	})

	recId := 0
	// apply calc function to each row of dataset
	df.Rapply(func(ser series.Series) series.Series {

		sphericalFraction := ser.Elem(aeronet.IdxSphrericalFract).Float()

		reIdx := ser.Subset(aeronet.IdxReM).Float()
		imIdx := ser.Subset(aeronet.IdxImM).Float()

		DateStr := ser.Subset(aeronet.IdxDate).String()
		TimeStr := ser.Subset(aeronet.IdxTime).String()
		DateStr = strings.ReplaceAll(DateStr, "[", "")
		DateStr = strings.ReplaceAll(DateStr, "]", "")
		DateStrItems := strings.Split(DateStr, ":")
		DateStr = DateStrItems[2] + "-" + DateStrItems[1] + "-" + DateStrItems[0]
		TimeStr = strings.ReplaceAll(TimeStr, "[", "")
		TimeStr = strings.ReplaceAll(TimeStr, "]", "")
		DateStr = DateStr + "T" + TimeStr + ".000Z"

		// Add code to parse data

		tmpdate, err := time.Parse("2006-01-02T15:04:05.000Z", DateStr)
		if err != nil {
			log.Println(err)
		}
		// Interpolate refractive index
		// Real part
		_ = pl.Fit(aeronet.Wvl, reIdx)
		a.dll.SetRn(pl.Predict(a.Wvl))

		// Imaginary part
		_ = pl.Fit(aeronet.Wvl, imIdx)
		a.dll.SetRk(pl.Predict(a.Wvl))

		a.dll.SetSd(ser.Subset(aeronet.IdxDvDlnr).Float())
		a.dll.DoCalc(a.dll.NDP())

		// извлекаем данные об оптической толще
		AerOptThickness := ser.Subset(aeronet.IdxAot).Float()

		// вычисляем десятичный логарифм длины волны и оптической толщи для аппроксимации
		// с учетом формулы Ангстрема
		for i := range AerOptThickness {
			LgWavelength[i] = math.Log10(aeronet.Wvl[i])
			LgAerOptThickness[i] = math.Log10(AerOptThickness[i])
		}

		//Добавляем результаты расчета в список
		tmp := a.dll.CalcResult()
		tmp.RecordId = recId
		tmp.Dt = tmpdate
		tmp.SphericalFraction = sphericalFraction
		SpheroidList.PushBack(tmp)

		recId++
		return ser
	})

	// Сбрасываем настройки и готовимся к загрузке сферических матриц
	a.dll.ResetToDefaultState()
	a.dll.ReadConfig(false)
	a.SetWl(a.Wvl)
	a.dll.AllocateMemory()

	recId = 0
	df.Rapply(func(ser series.Series) series.Series {

		sphericalFraction := ser.Elem(aeronet.IdxSphrericalFract).Float()

		reIdx := ser.Subset(aeronet.IdxReM).Float()
		imIdx := ser.Subset(aeronet.IdxImM).Float()

		// Interpolate refractive index
		// Real part
		_ = pl.Fit(aeronet.Wvl, reIdx)
		a.dll.SetRn(pl.Predict(a.Wvl))

		// Imaginary part
		_ = pl.Fit(aeronet.Wvl, imIdx)
		a.dll.SetRk(pl.Predict(a.Wvl))

		a.dll.SetSd(ser.Subset(aeronet.IdxDvDlnr).Float())
		a.dll.DoCalc(a.dll.NDP())

		// извлекаем данные об оптической толще
		AerOptThickness := ser.Subset(aeronet.IdxAot).Float()

		// вычисляем десятичный логарифм длины волны и оптической толщи для аппроксимации
		// с учетом формулы Ангстрема
		for i := range AerOptThickness {
			LgWavelength[i] = math.Log10(aeronet.Wvl[i])
			LgAerOptThickness[i] = math.Log10(AerOptThickness[i])
		}

		// Самое время сохранить наши матрицы
		// сохраняем в расширенном формате, то есть с нулевыми столбцами
		tmp := a.dll.CalcResult()
		tmp.RecordId = recId
		tmp.SphericalFraction = sphericalFraction
		SpheresList.PushBack(tmp)

		recId++
		return ser
	})

	fmt.Printf("Len1 = %d, Len2 = %d\n", SpheresList.Size(),
		SpheresList.Size())

	// Итак, мы имеем два списка с результатами моделирования
	// legacy.SpheroidsList и legacy.SpheresList
	// номер элемента в каждом из списков соответствует номеру измерений
	tmpSpheroid := SpheroidList.Front()
	tmpSphere := SpheresList.Front()

	// Iterate over lists and combine data
	for (tmpSphere != nil) && (tmpSpheroid != nil) {
		calcResSpheroid := tmpSpheroid.Value()
		calcResSphere := tmpSphere.Value()
		SphericalFraction := calcResSphere.SphericalFraction / 100.0

		if calcResSphere.RecordId == calcResSpheroid.RecordId {
			combCalcRes := &calcresult.CalculusResult{RecordId: calcResSphere.RecordId}

			combCalcRes.Ext = SphericalFraction*calcResSphere.Ext +
				(1.0-SphericalFraction)*calcResSpheroid.Ext

			combBsc := SphericalFraction*calcResSphere.Ext/calcResSphere.Lr +
				(1.0-SphericalFraction)*calcResSpheroid.Ext/calcResSpheroid.Lr

			combCalcRes.Sca = SphericalFraction*calcResSphere.Sca +
				(1.0-SphericalFraction)*calcResSpheroid.Sca

			combCalcRes.Absb = SphericalFraction*calcResSphere.Absb +
				(1.0-SphericalFraction)*calcResSpheroid.Absb

			combCalcRes.Lr = combCalcRes.Ext / combBsc
			var c1, c2, c3 mat.Dense
			c1.Scale(SphericalFraction, calcResSphere.MuellerMat)
			c2.Scale(1.0-SphericalFraction, calcResSpheroid.MuellerMat)
			c3.Add(&c1, &c2)
			combCalcRes.MuellerMat = &c3
			r, _ := combCalcRes.MuellerMat.Dims()
			combCalcRes.MuL = ((combCalcRes.MuellerMat.At(r-1, 0) - combCalcRes.MuellerMat.At(r-1, 2)) /
				(combCalcRes.MuellerMat.At(r-1, 0) + combCalcRes.MuellerMat.At(r-1, 2))) * 100

			combCalcRes.VolC = calcResSpheroid.VolC
			combCalcRes.Angle = calcResSpheroid.Angle
			combCalcRes.SphericalFraction = SphericalFraction
			combCalcRes.Dt = calcResSpheroid.Dt

			if DEBUG {
				if SphericalFraction < 0.1 {
					fmt.Printf("Aot = %5.2f, LR=%5.2f, MUL=%5.2f, LR=%5.2f, MUL=%5.2f\n",
						combCalcRes.Ext, combCalcRes.Lr, combCalcRes.MuL, calcResSpheroid.Lr, calcResSpheroid.MuL)
				}
			}

			CombList.PushBack(combCalcRes)
		}
		tmpSpheroid = tmpSpheroid.Next()
		tmpSphere = tmpSphere.Next()
	}

	fmt.Printf("Len3 = %d\n", CombList.Size())

	if err := a.SaveResults(CombList, matdir, picdir); err != nil {
		fmt.Println("Ошибка сохранения файлов")
	}
}

// SetWl прокси до внитреннего метода dll.
// Устанавливает длину волны
func (a *MuellerMatrixAERONET) SetWl(wvl float64) {

	if a.Wvl != wvl {
		if wvl < 0.440 {
			wvl = 0.440
		} else if wvl > 1.064 {
			wvl = 1.064
		}
		a.dll.SetWl(wvl)
	}
}

func (a *MuellerMatrixAERONET) Finalize() {
	a.dll.ClearNDP()
	a.dll.DeallocateMemory()
}

func (a *MuellerMatrixAERONET) SaveResults(lst *calcresultlist.CalcResultsList, matdir, picdir string) error {
	// Check for out dir
	if _, err := os.Stat(matdir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(matdir, 0755); err != nil {
				log.Fatal(err)
			}
		}
	}
	// Check for pic dir
	if _, err := os.Stat(picdir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(picdir, 0755); err != nil {
				log.Fatal(err)
			}
		}
	}
	err := lst.ApplyForward(func(cr *calcresult.CalculusResult) {
		fname := fmt.Sprintf("%s/%s_%04d%02d%02dT%02d%02d%02d.out", matdir, lst.Prefix, cr.Dt.Year(),
			cr.Dt.Month(), cr.Dt.Day(), cr.Dt.Hour(), cr.Dt.Minute(), cr.Dt.Second())
		saveto := fmt.Sprintf("%s/%s_%04d%02d%02dT%02d%02d%02d", picdir, lst.Prefix, cr.Dt.Year(),
			cr.Dt.Month(), cr.Dt.Day(), cr.Dt.Hour(), cr.Dt.Minute(), cr.Dt.Second())
		fout, err := os.Create(fname)

		if err != nil {
			fmt.Printf("Error creating file %s, %s\n", fname, err)
			return
		}

		defer func() {
			err := fout.Close()
			if err != nil {
				fmt.Printf("Error closing file %s, err=%s", fname, err)
			}
		}()

		angle := cr.Angle
		M := cr.MuellerMat
		rows, _ := M.Dims()
		Vc := cr.VolC
		_, _ = fmt.Fprintf(fout, "%9.3e\t%9.3e\t%9.3e\t%9.3e\t%9.3e\t%9.3e"+
			"# Sca/V, Ext/V, Absb/V, LR, MuL, V \n",
			cr.Sca/Vc, cr.Ext/Vc, cr.Absb/Vc, cr.Lr, cr.MuL, cr.VolC)
		_, _ = fmt.Fprintf(fout, "%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t"+
			"%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\n",
			"Angle", "S11", "S12", "S13", "S14", "S21", "S22", "S23",
			"S24", "S31", "S32", "S33", "S34", "S41", "S42", "S43",
			"S44")
		for i := 0; i < rows; i++ {
			_, _ = fmt.Fprintf(fout, "%9.3f\t", angle[i])
			_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 0)/Vc)  //S11
			_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, 1)/Vc)  //S12
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

		_ = cr.DoPlotPolarization(saveto)
	})
	return err
}

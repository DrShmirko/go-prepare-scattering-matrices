package scattlib

import (
	"bufio"
	"fmt"
	"github.com/kshmirko/prepare-mueller-matrices/calcresult"
	"github.com/kshmirko/prepare-mueller-matrices/calcresultlist"

	"gonum.org/v1/gonum/mat"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/kshmirko/prepare-mueller-matrices/aeronet"
	"github.com/kshmirko/prepare-mueller-matrices/legacy"
	"gonum.org/v1/gonum/interp"
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
func (a *MuellerMatrixAERONET) Run(fname string, sf float64, skiprows int) {
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

		// this is unused code
		//// Вычисляем аппроксимационный полином
		//pf.SetXY(LgWavelength, LgAerOptThickness)
		//_ = pf.Fit()
		//
		//// Выводим значения оптической толши расчетные и после аппроксимации
		//// Этот участок кода необходим для отладки
		//if DEBUG {
		//	fmt.Printf("AOT = %.2f, %.2f \n", a.dll.Xext(),
		//		math.Pow(10.0, pf.Evaluate(math.Log10(a.Wvl))))
		//}

		//Добавляем результаты расчета в список
		tmp := a.dll.CalcResult()
		tmp.RecordId = recId
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

		// this is unused code =====
		// Вычисляем аппроксимационный полином
		//pf.SetXY(LgWavelength, LgAerOptThickness)
		//_ = pf.Fit()
		//
		//// Выводим значения оптической толши расчетные и после аппроксимации
		//// Этот участок кода необходим для отладки
		//if DEBUG {
		//	fmt.Printf("AOT = %.2f, %.2f \n", a.dll.Xext(),
		//		math.Pow(10.0, pf.Evaluate(math.Log10(a.Wvl))))
		//}
		// ==========

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

	if err := CombList.SaveResults(); err != nil {
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

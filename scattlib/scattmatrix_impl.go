package scattlib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	//df := dataframe.ReadCSV(ioContent, dataframe.HasHeader(true))

	df1, err := aeronet.NewAeronetDatasets(fname, skiprows)
	if err != nil {
		log.Fatal(err)
	}

	// filter out values
	df1 = df1.Filter(func(v *aeronet.AeronetDataset) bool {
		return v.SphericalFraction < sf
	})

	// Prepare for interpolation of refractive index
	pl := interp.PiecewiseLinear{}

	// Initialize lists for storage processed items
	SpheroidList := calcresultlist.NewCalcResultsList("sphrd")
	SpheresList := calcresultlist.NewCalcResultsList("sphrs")
	CombList := calcresultlist.NewCalcResultsList("total")

	recId := 0

	df1.Apply(func(v aeronet.AeronetDataset) {

		_ = pl.Fit(aeronet.Wvl, v.ReIdx)
		a.dll.SetRn(pl.Predict(a.Wvl))

		// Imaginary part
		_ = pl.Fit(aeronet.Wvl, v.ImIdx)
		a.dll.SetRk(pl.Predict(a.Wvl))

		// Set Volume SD to library storage
		a.dll.SetSd(v.VolSd)
		a.dll.DoCalc(a.dll.NDP())

		//Добавляем результаты расчета в список
		tmp := a.dll.CalcResult()
		tmp.RecordId = recId
		tmp.Dt = v.TimePoint
		tmp.SphericalFraction = v.SphericalFraction
		SpheroidList.PushBack(tmp)

		recId++
	})

	// Сбрасываем настройки и готовимся к загрузке сферических матриц
	a.dll.ResetToDefaultState()
	a.dll.ReadConfig(false)
	a.SetWl(a.Wvl)
	a.dll.AllocateMemory()

	recId = 0

	df1.Apply(func(v aeronet.AeronetDataset) {

		// Interpolate refractive index
		// Real part
		_ = pl.Fit(aeronet.Wvl, v.ReIdx)
		a.dll.SetRn(pl.Predict(a.Wvl))

		// Imaginary part
		_ = pl.Fit(aeronet.Wvl, v.ImIdx)
		a.dll.SetRk(pl.Predict(a.Wvl))

		a.dll.SetSd(v.VolSd)
		a.dll.DoCalc(a.dll.NDP())

		// Самое время сохранить наши матрицы
		// сохраняем в расширенном формате, то есть с нулевыми столбцами
		tmp := a.dll.CalcResult()
		tmp.RecordId = recId
		tmp.SphericalFraction = v.SphericalFraction
		tmp.Dt = v.TimePoint
		SpheresList.PushBack(tmp)

		recId++
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

// SaveResults - save results to pic and mat files
// Pictures are generated in svg, png and txt format
// aerosol scattering file has following structure
//   * 1st line: SSA, EXT
//   * 2nd line: Headers
//   * 3rd line: number of angles
//   * Angle, S11, S12, S13, S14, S21, S22, S33, ..., S44
// SSA = SCA/EXT, ABS=EXT-SCA
func (a *MuellerMatrixAERONET) SaveResults(lst *calcresultlist.CalcResultsList, matdir, picdir string) error {
	// Check for out dir
	if _, err := os.Stat(matdir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(matdir, 0755); err != nil {
				return err
			}
		}
	}
	// Check for pic dir
	if _, err := os.Stat(picdir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(picdir, 0755); err != nil {
				return err
			}
		}
	}

	// Apply function to every item of list
	err := lst.ApplyForward(func(cr *calcresult.CalculusResult) {
		fname := fmt.Sprintf("%s/%s_%04d%02d%02dT%02d%02d%02d.out", matdir, lst.Prefix, cr.Dt.Year(),
			cr.Dt.Month(), cr.Dt.Day(), cr.Dt.Hour(), cr.Dt.Minute(), cr.Dt.Second())
		saveto := fmt.Sprintf("%s/%s_%04d%02d%02dT%02d%02d%02d", picdir, lst.Prefix, cr.Dt.Year(),
			cr.Dt.Month(), cr.Dt.Day(), cr.Dt.Hour(), cr.Dt.Minute(), cr.Dt.Second())

		if err := cr.DoSaveMatrixToFile(fname); err != nil {
			log.Fatalf("Error %s", err)
		}

		if err := cr.DoPlotPolarization(saveto); err != nil {
			log.Fatalf("Error %s", err)
		}

	})
	return err
}

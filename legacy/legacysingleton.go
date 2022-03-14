package legacy

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"unsafe"

	"github.com/kshmirko/prepare-mueller-matrices/calcresult"

	"gonum.org/v1/gonum/mat"
)

//#cgo LDFLAGS: -L ../lib -lspheroid
//#cgo CFLAGS: -I ../include/
//#include "legacy.h"
//#include <stdlib.h>
import "C"

// carray2slice Отображает массив типа float из  С на slice в  Го
// при этом слайс и массив использут один и тот же кусок памяти
// по сути, мы конструируем слайс используя информацию о размере
// массива, его емкости и указателе на память, в котрой хранится
// массив. В этом подходе мы должны озаботиться, чтобы сборщик
// мусора не уничтожил наши данные (указатель Data).
// Но поскольку мы рассматриывем отображение C массива, а он не
// обрабатывается сборщиком мусора, здесь волноваться не приходится.
func carray2slice(array *C.float, len int) ([]float32, error) {

	var list []float32

	if array == nil {
		return nil, errors.New("cи массив - указатель на NULL")
	}
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&list))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(array))
	return list, nil
}

// float32float64 Преобразует слайс типа float32 в слайт типа float64
// здесь выполняется копирование с приведением типа кажждого элемента к float64
func float32float64(array []float32) []float64 {
	N := len(array)
	list := make([]float64, N)

	for i := 0; i < N; i++ {
		list[i] = float64(array[i])
	}
	return list
}

// Instance возвращает сущность объекта, который создается один раз
// и существует в программе в единственном экземпляре
func Instance() *Singleton {
	once.Do(func() {
		instance = new(Singleton)
		instance.IsAllocated = false
	})
	return instance
}

// ReadConfigFile читает файл конфигурации инициализируя внутренние
// переменные модуля. Чтение идет по вызову внутренней процедуры
// модуля libspheroid (функция dls_read_config_)
func (s *Singleton) ReadConfigFile(fname string) {
	cfname := C.CString(fname)
	N := len(fname)
	cN := C.int(N)
	defer C.free(unsafe.Pointer(cfname))
	C.dls_read_config_(cfname, &cN)
}

func (s *Singleton) SetDistnameO(fname string) {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	C.set_distname_o(cfname)
}

func (s *Singleton) SetDistnameF(fname string) {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	C.set_distname_f(cfname)
}

func (s *Singleton) SetDistnameN(fname string) {
	cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(cfname))
	C.set_distname_n(cfname)
}

func (s *Singleton) SetCommName(id int, data string) {
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	tmp := C.int(id)
	C.set_comm_name(&tmp, cdata)
}

// AllocateMemory выделение памяти под внутренние переменные в куче
// вызывать строго после ReadInitFile
func (s *Singleton) AllocateMemory() {
	if !s.IsAllocated {
		KA := C.int(1)
		C.alloc_dls_array_(&C.key,
			&C.keyel,
			&KA)
		s.IsAllocated = true
	}
}

// DeallocateMemory освобождение памяти вынутренних объектов
// вызывать после AllocateMemory в конце программы
func (s *Singleton) DeallocateMemory() {
	if s.IsAllocated {
		KA := C.int(2)
		C.alloc_dls_array_(&C.key,
			&C.keyel,
			&KA)
		s.IsAllocated = false
	}
}

// ResetToDefaultState - сбрасывает настройки солвера в начальное состояние
// чтобы далее его использовать, нужно будет снова прочитать конфигурацию и
// инициализиоровать переменные.
func (s *Singleton) ResetToDefaultState() {
	s.ClearNDP()
	s.DeallocateMemory()
}

// DoCalc Perform calculations
// if ndp = 0, loads database from disk before
// else if ndp=1, then data already loaded into memory
func (s *Singleton) DoCalc(ndp int) int {
	if !s.IsAllocated {
		fmt.Printf("Внутренние переменые сперва должны быть инициализированы")
	}
	C.ndp = C.int(ndp)
	if ndp == 0 {
		log.Println("\nЗагружаем ядра с диска")
	}
	C.optchar_(&C.ndp)
	return int(C.ndp)
}

// RatiosX - возвращает отношение полуосей элипсоида
func (s *Singleton) RatiosX() ([]float64, error) {
	vector, err := carray2slice(&C.r[0], int(C.kr))

	if err != nil {
		return []float64{}, err
	}
	return float32float64(vector), nil
}

// RatiosY - возвращает относительный вклад в концентрацию
func (s *Singleton) RatiosY() ([]float64, error) {
	vector, err := carray2slice(&C.rd[0], int(C.kr))

	if err != nil {
		return []float64{}, err
	}
	return float32float64(vector), nil
}

// Sd - значения ординат функции распределения, заданные таблично
func (s *Singleton) Sd() ([]float64, error) {
	vector, err := carray2slice(&C.sd[0], int(C.kn))

	if err != nil {
		return []float64{}, err
	}
	return float32float64(vector), nil
}

// Grid - значение абсцисс функции распределения, заданные таблично
func (s *Singleton) Grid() ([]float64, error) {
	vector, err := carray2slice(&C.grid[0], int(C.kn))

	if err != nil {
		return []float64{}, err
	}
	return float32float64(vector), nil
}

// SetPSD -
func (s *Singleton) SetPSD(x, y []float64) {
	N := len(x)
	for i := 0; i < N; i++ {
		C.grid[i] = C.float(x[i])
		C.sd[i] = C.float(y[i])
	}

	C.kn = C.int(N)
}

// SetSd -
func (s *Singleton) SetSd(y []float64) {
	N := len(y)
	for i, v := range y {
		C.sd[i] = C.float(v)
	}
	C.kn = C.int(N)
}

// Rrr Значение абсцисс функции распределения, в случае если
// функция задается как логнормальное распределение
func (s *Singleton) Rrr() ([]float64, error) {
	vector, err := carray2slice(&C.rrr[0], int(C.kn))

	if err != nil {
		return []float64{}, err
	}

	return float32float64(vector), nil
}

// Ar Значение ординат функции распределения, в случае если
// функция задается как логнормальное распределение
func (s *Singleton) Ar() ([]float64, error) {
	vector, err := carray2slice(&C.ar[0], int(C.kn))

	if err != nil {
		return []float64{}, err
	}

	return float32float64(vector), nil
}

// Xgrid Значение отсчетов, похоже это значения, на что интерполируется Sd
func (s *Singleton) Xgrid() ([]float64, error) {
	vector, err := carray2slice(&C.xgrid[0], int(C.kn))

	if err != nil {
		return []float64{}, err
	}
	return float32float64(vector), nil
}

// Angle возвращает вектор углов рассеяния, для которых выполнялись расчеты
// матрицы
func (s *Singleton) Angle() ([]float64, error) {
	vector, err := carray2slice(&C.angle[0], int(C.km))

	if err != nil {
		return []float64{}, err
	}

	return float32float64(vector), nil
}

/**
F11, F12, F22, F33, F44, F34 - ненулевые элементы матрицы рассеяния,
нормированные на F11
*/

// F11 -
func (s *Singleton) F11() []float64 {
	vector, err := carray2slice(&C.f11[0], int(C.km))

	if err != nil {
		return []float64{}
	}

	return float32float64(vector)
}

// F12 -
func (s *Singleton) F12() []float64 {
	vector, err := carray2slice(&C.f12[0], int(C.km))

	if err != nil {
		return []float64{}
	}

	return float32float64(vector)
}

// F22 -
func (s *Singleton) F22() []float64 {
	vector, err := carray2slice(&C.f22[0], int(C.km))

	if err != nil {
		return []float64{}
	}

	return float32float64(vector)
}

// F33 -
func (s *Singleton) F33() []float64 {
	vector, err := carray2slice(&C.f33[0], int(C.km))

	if err != nil {
		return []float64{}
	}

	return float32float64(vector)
}

// F34 -
func (s *Singleton) F34() []float64 {
	vector, err := carray2slice(&C.f34[0], int(C.km))

	if err != nil {
		return []float64{}
	}

	return float32float64(vector)
}

// F44 -
func (s *Singleton) F44() []float64 {
	vector, err := carray2slice(&C.f44[0], int(C.km))

	if err != nil {
		return []float64{}
	}

	return float32float64(vector)
}

// Ac В случае, если функция распределения задавалась при помощи формул
// вернет значение концентрации
func (s *Singleton) Ac() float64 {
	return float64(C.ac)
}

// AcAeronet В случае, если функция распределения задавалась таблично и отсчеты
// соответствуют отсчетам в файле aeronet, то функция вернет значение
// объемной концентрации
func (s *Singleton) AcAeronet() float64 {
	dlnr := 0.2716
	ret := 0.0
	y, _ := s.Sd()

	for _, yi := range y {
		ret += yi
	}
	ret *= dlnr
	return ret
}

// Ac0 Вычисляет объемную концентрацию
// требуется, чтобы таблично заданная функция была вида dV/dlnr
func (s *Singleton) Ac0() float64 {
	x, _ := s.Grid()
	y, _ := s.Sd()
	N := len(x)
	ret := 0.0

	for i := 1; i < N; i++ {
		dlnr := math.Log(x[i]) - math.Log(x[i-1])
		Ym := 0.5 * (y[i] + y[i-1])
		ret += Ym * dlnr
	}
	return ret
}

// Cm возвращает стлайс с концентрациями для каждой из используемых мод
func (s *Singleton) Cm() ([]float64, error) {
	vector, err := carray2slice(&C.cm[0], int(C.nmd))

	if err != nil {
		return []float64{}, err
	}

	return float32float64(vector), nil
}

// GetSm возвращает стлайс с полуширинами для каждой из используемых мод
func (s *Singleton) GetSm() ([]float64, error) {
	vector, err := carray2slice(&C.sm[0], int(C.nmd))

	if err != nil {
		return []float64{}, err
	}

	return float32float64(vector), nil
}

// Rmm возвращает стлайс с медианными радиусами для каждой из используемых мод
func (s *Singleton) Rmm() ([]float64, error) {
	vector, err := carray2slice(&C.rmm[0], int(C.nmd))

	if err != nil {
		return []float64{}, err
	}

	return float32float64(vector), nil
}

// SetLNPars аналогичная функция, но на вход подаются слайсы каждого из параметров
// распределений
func (s *Singleton) SetLNPars(cm, sm, rmm []float64) {
	N := len(cm)
	if N > 2 {
		N = 2
	}

	for i := 0; i < N; i++ {
		C.cm[i] = C.float(cm[i])
		C.sm[i] = C.float(sm[i])
		C.rmm[i] = C.float(rmm[i])
	}

	C.nmd = C.int(N)
}

// Xext  Extinction crosssection mkm^2
func (s *Singleton) Xext() float64 {
	return float64(C.xext)
}

// Xabs Absorbtion crossection mkm^2
func (s *Singleton) Xabs() float64 {
	return float64(C.xabs)
}

// Xsca Scattering crosssection mkm^2
func (s *Singleton) Xsca() float64 {
	return float64(C.xsca)
}

// Ssa Single scattering albedo
func (s *Singleton) Ssa() float64 {
	return float64(C.albedo)
}

// Lr Lidar Ratio
func (s *Singleton) Lr() float64 {
	return float64(C.xblr)
}

// Ldr Linear depolarization ratio
func (s *Singleton) Ldr() float64 {
	return float64(C.xldr)
}

// Wl Возвращает длину волны мкм
func (s *Singleton) Wl() float64 {
	return float64(C.wl)
}

// SetWl Устанавливает длину волны в мкм
func (s *Singleton) SetWl(wl float64) {
	C.wl = C.float(wl)
}

// RefrIdx  return complex refractive index
func (s *Singleton) RefrIdx() complex128 {
	return complex(s.Rn(), s.Rk())
}

// SetRefrIdx set rn and rk
func (s *Singleton) SetRefrIdx(midx complex128) {
	s.SetRn(real(midx))
	s.SetRk(imag(midx))
}

// Rn Возвращает и устанавливает действительную часть показателя преломления
func (s *Singleton) Rn() float64 {
	return float64(C.rn)
}

// SetRn -
func (s *Singleton) SetRn(rn float64) {
	C.rn = C.float(rn)
}

// Rk Возвращает и устанавливает мнимую часть показателя преломления
func (s *Singleton) Rk() float64 {
	return float64(C.rk)
}

// SetRk -
func (s *Singleton) SetRk(rk float64) {
	C.rk = C.float(rk)
}

// NDP - возвращает параметр, контроллирующий загрузку базы данных
func (s *Singleton) NDP() int {
	return int(C.ndp)
}

// IsSetNDP -  test weather ndp is set
func (s *Singleton) IsSetNDP() bool {
	return int(C.ndp) == 1
}

// SetNDP - Устанавливает параметр, контроллирующий загрузку базы данных
func (s *Singleton) SetNDP() {
	C.ndp = C.int(1)
}

// ClearNDP - set ndp to zero
func (s *Singleton) ClearNDP() {
	C.ndp = C.int(0)
}

// PrintKeyParams вывод настроек расчета функций
func (s *Singleton) PrintKeyParams() {
	fmt.Println("Key parameters:=============")
	fmt.Printf("key = %d\nkey_rd = %d\nkeyel = %d\nkeysub = %d\nkeyls = %d\nkey_org = %d\n",
		int(C.key), int(C.key_rd), int(C.keyel), int(C.keysub), int(C.keyls), int(C.key_org))
	fmt.Printf("key_fx = %d\nkey_grid1 = %d\nkey_rd1 = %d\nkn = %d\nkm = %d\nkr = %d\nnratn = %d\n\n",
		int(C.key_fx), int(C.key_grid1), int(C.key_rd1), int(C.kn), int(C.km), int(C.kr), int(C.nratn))

}

// MuellerMatrix / Возвращает угол рассеяния и матрицу рассеяния
func (s *Singleton) MuellerMatrix() ([]float64, *mat.Dense) {
	ret := mat.NewDense(int(C.km), 16, nil)

	angle, _ := s.Angle()

	f11 := s.F11()
	f12 := s.F12()
	f22 := s.F22()
	f33 := s.F33()
	f34 := s.F34()
	f44 := s.F44()
	N := len(angle)
	f43 := make([]float64, N)

	// умножаем вектор на раскрываем нашу фазовую матрицу в матрицу
	// рассеяния, при этом интеграл первого элемента ее по телесному
	// углу должен быть равен АОТ
	for i := 0; i < N; i++ {
		f12[i] = -f12[i] * f11[i]
		f22[i] = f22[i] * f11[i]
		f33[i] = f33[i] * f11[i]
		f34[i] = f34[i] * f11[i]
		f43[i] = -f34[i]
		f44[i] = f44[i] * f11[i]
	}

	// заполняем ненулевые столбцы матрицы Мюллера значениями
	ret.SetCol(0, f11)
	ret.SetCol(1, f12)
	ret.SetCol(4, f12)
	ret.SetCol(5, f22)
	ret.SetCol(10, f33)
	ret.SetCol(11, f34)
	ret.SetCol(14, f43)
	ret.SetCol(15, f44)

	return angle, ret
}

func (s *Singleton) MuellerMatrixShort() ([]float64, *mat.Dense) {
	ret := mat.NewDense(int(C.km), 6, nil)

	angle, _ := s.Angle()
	//AOT := s.Xsca()
	f11 := s.F11()
	f12 := s.F12()
	f22 := s.F22()
	f33 := s.F33()
	f34 := s.F34()
	f44 := s.F44()
	N := len(angle)

	// умножаем вектор на раскрываем нашу фазовую матрицу в матрицу
	// рассеяния, при этом интеграл первого элемента ее по телесному
	// углу должен быть равен АОТ
	for i := 0; i < N; i++ {
		f12[i] = -f12[i] * f11[i]
		f22[i] = f22[i] * f11[i]
		f33[i] = f33[i] * f11[i]
		f34[i] = f34[i] * f11[i]
		f44[i] = f44[i] * f11[i]
	}

	ret.SetCol(0, f11)
	ret.SetCol(1, f12)
	ret.SetCol(2, f22)
	ret.SetCol(3, f33)
	ret.SetCol(4, f34)
	ret.SetCol(5, f44)

	return angle, ret
}

// SaveResult сохранаяе результаты расчета в файл для дальнейшей постобработки
func (s *Singleton) SaveResult(prefix string, idx int) {
	fname := fmt.Sprintf("%s_%05d.dat", prefix, idx)
	fout, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := fout.Close()
		if err != nil {
			fmt.Printf("Error closing file")
			panic(err)
		}
	}()

	angle, M := s.MuellerMatrix()
	rows, cols := M.Dims()
	Vc := s.Ac0()
	_, _ = fmt.Fprintf(fout, "%9.3e\t%9.3e\t%9.3e\t%9.3e\t%9.3e\t# \n",
		s.Xext()/Vc, s.Xsca()/Vc, s.Xabs()/Vc, s.Lr(), s.Ldr())
	_, _ = fmt.Fprintf(fout, "%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t"+
		"%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\t%9s\n",
		"Angle", "S11", "S12", "S13", "S14", "S21", "S22", "S23",
		"S24", "S31", "S32", "S33", "S34", "S41", "S42", "S43",
		"S44")
	for i := 0; i < rows; i++ {
		_, _ = fmt.Fprintf(fout, "%9.3f\t", angle[i])
		for j := 0; j < cols; j++ {
			_, _ = fmt.Fprintf(fout, "%9.3e\t", M.At(i, j)/Vc)
		}
		_, _ = fmt.Fprintf(fout, "\n")
	}
}

func (s *Singleton) CalcResult() *calcresult.CalculusResult {
	vc := s.Ac0()
	angle, mm := s.MuellerMatrixShort()
	ret := calcresult.NewCalcResult(false,
		0, s.Xext(), s.Xsca(), s.Xabs(), vc, s.Lr(), s.Ldr(), mm, angle, 1)

	return ret
}

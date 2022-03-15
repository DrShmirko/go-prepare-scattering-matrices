package aeronet

import (
	"bufio"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// NewAeronetDatasets - читает из файла данные расчетов AERONET. Пропускает skiprows
// строк. Возвращает масив AeronetDatasets
func NewAeronetDatasets(filename string, skiprows int) (AeronetDatasets, error) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ioreader := bufio.NewReader(strings.NewReader(string(content)))

	res := make(AeronetDatasets, 0, 20)

	lineno := 0
	for {
		line, err := ioreader.ReadString('\n')
		if err != nil {
			break
		}
		lineno++

		if lineno > skiprows {
			items := strings.Split(line, ",")

			tmp := AeronetDataset{
				TimePoint:         extract_timepoint(items),
				AotExt:            extract_field(IdxAot, items),
				ReIdx:             extract_field(IdxReM, items),
				ImIdx:             extract_field(IdxImM, items),
				VolSd:             extract_field(IdxDvDlnr, items),
				SphericalFraction: extract_sphericalfraction(items),
			}

			res = append(res, tmp)

		}
	}

	return res, nil
}

// extract_timepoint - из списка токенов выделяет дату и время и
// объединет в одну переменную
func extract_timepoint(items []string) time.Time {

	DateStrItems := strings.Split(items[IdxDate], ":")
	dateStr := DateStrItems[2] + "-" + DateStrItems[1] + "-" + DateStrItems[0]
	dateStr = dateStr + "T" + items[IdxTime] + ".000Z"

	ret, err := time.Parse("2006-01-02T15:04:05.000Z", dateStr)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

// extract_sphericalfraction - извлекает сферичность из набора токенов
func extract_sphericalfraction(items []string) float64 {

	ret, err := strconv.ParseFloat(items[IdxSphrericalFract], 64)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

// extract_fields - извлекае поле, указанное idx из набора лексем
func extract_field(idx []int, items []string) []float64 {
	ret := make([]float64, len(idx))
	var err error

	for i, v := range idx {
		if ret[i], err = strconv.ParseFloat(items[v], 64); err != nil {
			log.Fatal(err)
		}
	}
	return ret
}

// Filter filter function - apply boolean function fo each element of the array AeronetDatasets
func (a *AeronetDatasets) Filter(fun func(v AeronetDataset) bool) AeronetDatasets {
	ret := make(AeronetDatasets, 0, 20)

	for _, v := range *a {
		if fun(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Apply filter function
func (a *AeronetDatasets) Apply(fun func(v *AeronetDataset)) {
	for i := 0; i < a.Size(); i++ {
		fun(&(*a)[i])
	}
}

// Apply filter function
func (a AeronetDatasets) ApplyCopy(fun func(v AeronetDataset) AeronetDataset) AeronetDatasets {
	for i := 0; i < a.Size(); i++ {
		a[i] = fun(a[i])
	}
	return a
}

func (a *AeronetDatasets) Size() int {
	return len(*a)
}

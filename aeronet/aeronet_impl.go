package aeronet

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func NewAeronetDatasets(filename string, skiprows int) (*AeronetDatasets, error) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ioreader := bufio.NewReader(strings.NewReader(string(content)))
	var res AeronetDatasets

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
				AotExt:            extract_fields(IdxAot, items),
				ReIdx:             extract_fields(IdxReM, items),
				ImIdx:             extract_fields(IdxImM, items),
				VolSd:             extract_fields(IdxDvDlnr, items),
				SphericalFraction: extract_sphericalfraction(items),
			}

			res = append(res, tmp)

		}
	}
	fmt.Println(len(res))
	return &res, nil
}

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

func extract_sphericalfraction(items []string) float64 {

	ret, err := strconv.ParseFloat(items[IdxSphrericalFract], 64)
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

func extract_fields(idx []int, items []string) []float64 {
	ret := make([]float64, len(idx))
	var err error

	for i, v := range idx {
		if ret[i], err = strconv.ParseFloat(items[v], 64); err != nil {
			log.Fatal(err)
		}
	}
	return ret
}

// Filter filter function
func (a *AeronetDatasets) Filter(fun func(v *AeronetDataset) bool) *AeronetDatasets {
	var ret AeronetDatasets

	for _, v := range *a {
		if fun(&v) {
			ret = append(ret, v)
		}
	}
	return &ret
}

// Apply filter function
func (a *AeronetDatasets) Apply(fun func(v AeronetDataset)) {
	for i := 0; i < a.Size(); i++ {
		fun((*a)[i])
	}
}

func (a *AeronetDatasets) Size() int {
	return len(*a)
}

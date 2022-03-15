package agglomerated

import (
	"fmt"

	"github.com/batchatco/go-native-netcdf/netcdf/cdf"
)

func LoadFromFile(cdfname string) (*SingleParticle, error) {
	var sp SingleParticle

	fin, err := cdf.Open(cdfname)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s", err)
	}

	defer fin.Close()

	/// Start reading ReIdx from netcdf
	mre, err := fin.GetVariable("mre")
	if err != nil {
		return nil, fmt.Errorf("there is no mre variable in netcdf %s", err)
	}

	scale_factor, has := mre.Attributes.Get("scale_factor")
	if !has {
		return nil, fmt.Errorf("there is no scale_factor attribute in mre")
	}

	scale_factor_val, ok := scale_factor.(float64)
	if !ok {
		return nil, fmt.Errorf("error casting scale_factor of mre value to float64")
	}

	mreval, ok := mre.Values.([]int16)

	if !ok {
		return nil, fmt.Errorf("error casting mre values to array of int16s")
	}

	// fill inernal array with loaded ones
	sp.ReIdx = make([]float64, len(mreval))
	for i := 0; i < len(mreval); i++ {
		sp.ReIdx[i] = float64(mreval[i]) * scale_factor_val
	}

	// read mim
	mim, err := fin.GetVariable("mim")
	if err != nil {
		return nil, fmt.Errorf("there is no mim variable in netcdf %s", err)
	}

	scale_factor, has = mim.Attributes.Get("scale_factor")
	if !has {
		return nil, fmt.Errorf("there is no scale_factor attribute in mre")
	}

	scale_factor_val, ok = scale_factor.(float64)
	if !ok {
		return nil, fmt.Errorf("error casting scale_factor of mim value to float64")
	}

	mimval, ok := mim.Values.([]int16)

	if !ok {
		return nil, fmt.Errorf("error casting mre values to array of int16s")
	}

	sp.ImIdx = make([]float64, len(mimval))
	for i := 0; i < len(mimval); i++ {
		sp.ImIdx[i] = float64(mimval[i]) * scale_factor_val
	}

	return &sp, nil
}

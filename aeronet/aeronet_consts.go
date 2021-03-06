package aeronet

const (
	SIZE_AOT = 4
	SIZE_SD  = 22
)

var (
	Wvl = []float64{0.440, 0.675, 0.870, 1.02}

	IdxDate            = 1
	IdxTime            = 2
	IdxDayOfYear       = 4
	IdxAot             = []int{10, 11, 12, 13}
	IdxReM             = []int{32, 33, 34, 35}
	IdxImM             = []int{36, 37, 38, 39}
	IdxSphrericalFract = 52
	IdxDvDlnr          = []int{53, 54, 55, 56, 57, 58, 59, 60, 61, 62,
		63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74}
)

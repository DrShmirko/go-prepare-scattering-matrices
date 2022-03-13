package main

import (
	"flag"
	"log"

	"github.com/kshmirko/prepare-mueller-matrices/scattlib"
)

func main() {
	sphericalFraction := flag.Float64("sf", 1,
		"Пороговое значение сферичности.")
	skipRows := flag.Int("skip", 6, "Сколько строк пропускать.")
	waveLen := flag.Float64("wvl", 0.870, "Длина волны для расчета матриц.")
	picDir := flag.String("picdir", "./pic", "Каталог, для сохранения графиков.")
	matDir := flag.String("matdir", "./out", "Каталог, для сохранения матриц.")

	flag.Parse()

	// ================ Testers and checkers =================

	// Test weather AERONET
	// file is presented in the command line
	if flag.NArg() != 1 {
		log.Println("Необхоимо указать хотя бы одно имя файла с данными.")
		flag.PrintDefaults()
		return
	}

	// Check weather spherical fraction in valid range
	if *sphericalFraction <= 0 || *sphericalFraction >= 100 {
		log.Println("Значение сферичности должно быть на отрезке [1.0; 99.0]")
		return
	}

	// Check weather waveLen in valid range
	if *waveLen < 0.440 || *waveLen > 1.020 {
		log.Println("Длина волны должна принадлежать отрезку [0.44; 1.02] мкм.")
		return
	}

	// check weather skipRows in valid range
	if *skipRows < 0 || *skipRows > 10 {
		log.Println("Число пропускаемых строк должно быть " +
			"положительно и меньше 10.")
	}

	// Initialization of solvers
	//_, err := agglomerated.LoadFromFile("mod_agglomerated_debris_shape_01.npz.cdf")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//_ = *picDir
	//_ = *matDir
	mc := scattlib.NewMuellerMatrixAERONET(*waveLen)
	mc.Run(flag.Arg(0), *sphericalFraction, *skipRows, *matDir, *picDir)
	mc.Finalize()

}

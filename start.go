package main

import (
	"flag"
	//"fmt"
	"github.com/kshmirko/prepare-mueller-matrices/scattlib"
	//num "gonum.org/v1/gonum"
	"log"
)

func main() {
	//version, _ := num.Version()
	//fmt.Printf("Версия GoNum: %s\n\n\n", version)

	// Define command line flags
	sphericalFraction := flag.Float64("sf", 1, "Установить пороговое значение сферичности")
	//isSpheres := flag.Bool("spheres", false, "Использовать сферические ядра или сфероидальные")
	skipRows := flag.Int("skip", 6, "Сколько строк пропускать")
	waveLen := flag.Float64("wvl", 0.870, "Длина волны для расчета матриц")

	flag.Parse()

	// ================ Testers and checkers =================

	// Test weather AERONET
	// file is presented in the command line
	if flag.NArg() != 1 {
		log.Println("Необхоимо указать хотя бы одно имя файла с данными.")
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
		log.Println("Число пропускаемых строк должно быть положительно и меньше 10.")
	}

	// Initialization of solvers

	mc := scattlib.NewMuellerMatrixAERONET(*waveLen)
	mc.Run(flag.Arg(0), *sphericalFraction, *skipRows)
	mc.Finalize()

}

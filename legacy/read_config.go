package legacy

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
)

//#cgo LDFLAGS: -L ../lib -lspheroid
//#cgo CFLAGS: -I ../include
//#include "legacy.h"
//#include <stdlib.h>
import "C"

//go:embed input_sphrs.dat
var spheresConfig string

//go:embed input_sphrds.dat
var spheroidsConfig string

// ReadConfig Читаем файл с настройками, который мы встраиваем в исполняемый файл
func (s *Singleton) ReadConfig(spheres bool) {
	var fin io.Reader

	// Check which config file to read
	if spheres {
		fin = strings.NewReader(spheresConfig)
	} else {
		fin = strings.NewReader(spheroidsConfig)
	}

	var (
		key, keyEL, keySUB, keyLS, key_org,
		key_fx, key_RD1, key_SD, id, nmd, kn,
		kr, km, nratn int

		wl, rn, rk, rgmin, rgmax, wlmin, wlmax, grid_i,
		sd_i, cm_i, sm_i, rmm_i, r_i, rd_i, ang_i float64
	)

	scanner := bufio.NewScanner(fin)

	scanner.Scan()
	str := scanner.Text()
	_, err := fmt.Sscanf(str, "%d%d%d%d%d%d%d",
		&key, &keyEL, &keySUB, &keyLS, &key_org, &key_fx, &key_RD1)

	if err != nil {
		log.Fatal("Error reading first line in config.")
	}

	scanner.Scan()
	str = scanner.Text()
	pos, err := fmt.Sscanf(str, "%f%f%f%f%f%f%f",
		&wl, &rn, &rk, &rgmin, &rgmax, &wlmin, &wlmax)

	if err != nil {
		log.Println(pos)
		log.Fatal("Error reading second line in config.")
	}

	scanner.Scan()
	str = scanner.Text()
	_, err = fmt.Sscanf(str, "%d %d %d",
		&key_SD, &id, &nmd)
	if err != nil {
		log.Fatalln("Error reading keySD, id, nmd")
	}
	C.nmd = C.int(nmd)
	C.id = C.int(id)
	C.key_sd = C.int(key_SD)
	if key_SD == 0 {
		scanner.Scan()
		scanner.Scan()
		str := scanner.Text()
		_, err := fmt.Sscanf(str, "%d", &kn)
		if err != nil {
			panic(err)
		}
		log.Println(kn)
		C.kn = C.int(kn)
		if kn > C.KNpar {
			log.Fatal("kn from input file should be less than KNpar")
		}
		for i := 0; i < kn; i++ {
			scanner.Scan()
			str := scanner.Text()
			_, err := fmt.Sscanf(str, "%f%f", &grid_i, &sd_i)
			if err != nil {
				panic(err)
			}

			C.sd[i] = C.float(sd_i)
			C.grid[i] = C.float(grid_i)
		}
	} else {
		log.Println(nmd)
		for i := 0; i < nmd; i++ {
			scanner.Scan()
			str := scanner.Text()
			_, err := fmt.Sscanf(str, "%f%f%f", &cm_i, &sm_i, &rmm_i)

			if err != nil {
				log.Fatalln(err)
			}

			C.cm[i] = C.float(cm_i)
			C.sm[i] = C.float(sm_i)
			C.rmm[i] = C.float(rmm_i)
		}
		scanner.Scan()
		str := scanner.Text()
		_, err := fmt.Sscanf(str, "%d", &kn)
		if err != nil {
			log.Fatalln(err)
		}

		C.kn = C.int(kn)
		for i := 0; i < kn; i++ {
			scanner.Scan()
			str := scanner.Text()
			_, err := fmt.Sscanf(str, "%f", &grid_i)
			if err != nil {
				log.Fatalln(err)
			}

			C.grid[i] = C.float(grid_i)
		}
	}
	for i := 0; i < kn; i++ {
		C.xgrid[i] = C.grid[i]
	}

	scanner.Scan()
	str = scanner.Text()
	log.Println(str)
	s.SetDistnameO(str)
	scanner.Scan()
	str = scanner.Text()
	log.Println(str)
	s.SetDistnameF(str)
	scanner.Scan()
	str = scanner.Text()
	log.Println(str)
	s.SetDistnameN(str)

	scanner.Scan()
	str = scanner.Text()
	_, err = fmt.Sscanf(str, "%d", &kr)
	if err != nil {
		panic(err)
	}

	C.kr = C.int(kr)
	if kr > C.KRpar {
		log.Fatal("kr from input file should be less than KRpar")
	}
	for i := 0; i < kr; i++ {
		scanner.Scan()
		str := scanner.Text()
		_, err := fmt.Sscanf(str, "%f%f", &r_i, &rd_i)
		if err != nil {
			panic(err)
		}

		C.r[i] = C.float(r_i)
		C.rd[i] = C.float(rd_i)
	}
	log.Println("R, Rd has been read")
	scanner.Scan()
	str = scanner.Text()
	_, err = fmt.Sscanf(str, "%d", &km)
	if err != nil {
		log.Fatal(err)
	}

	C.km = C.int(km)
	if km > C.KMpar {
		log.Fatal("km from input file should be less than KMpar")
	}
	for i := 0; i < km; i++ {
		scanner.Scan()
		str := scanner.Text()
		_, err := fmt.Sscanf(str, "%f", &ang_i)
		if err != nil {
			log.Fatal(err)
		}

		C.angle[i] = C.float(ang_i)
	}
	log.Println("Angles has been read")
	scanner.Scan()
	str = scanner.Text()
	_, err = fmt.Sscanf(str, "%d", &nratn)
	if err != nil {
		log.Fatal(err)
	}

	C.nratn = C.int(nratn)
	if nratn > C.KR1par {
		log.Fatal("nratn from input file should be less than KR1par")
	}

	for i := 0; i < nratn; i++ {
		scanner.Scan()
		str = scanner.Text()
		s.SetCommName(i, str)
	}
	log.Println("COMM_NAMES has been read")

	if (key == 4) && (key_org == 1) {
		log.Fatal("If key=4, key_org should be 0")
	}
	if (key_org == 1) && (key_fx == 1) {
		log.Fatal("STOP: key_org=1 & key_fx=1\n" +
			"If you want to use key_fx=1 change key_org to 0\n" +
			"If key_org=0 is your choice," +
			"check dir_name1 in matrix_fixget" +
			"and run the code again.\n")
	}

	if key == 2 {
		C.key_grid1 = C.int(1)
	} else {
		C.key_grid1 = C.int(0)
	}

	C.key = C.int(key)
	C.keyel = C.int(keyEL)
	C.keysub = C.int(keySUB)
	C.keyls = C.int(keyLS)
	C.key_org = C.int(key_org)
	C.key_fx = C.int(key_fx)
	C.key_rd1 = C.int(key_RD1)

	C.key_rd = C.int(1)

	C.wl = C.float(wl)
	C.rn = C.float(rn)
	C.rk = C.float(rk)
	C.pomin = C.float(math.Pi * rgmin / wlmax)
	C.pomax = C.float(math.Pi * rgmax / wlmin)

	C.key_sd = C.int(key_SD)
	C.id = C.int(id)
	C.nmd = C.int(nmd)

	if key_SD == 1 {
		var kn, ia, id, knpar, idbg C.int
		//f7 := os.NewFile(7, "LNPAR.dat")
		//f9 := os.NewFile(9, "SizeDis.dat")

		kn = -C.kn
		ia = C.int(1)
		id = C.id
		knpar = C.KNpar
		idbg = C.int(1)
		C.sizedisdn_(&kn, &ia, &id, &C.nmd,
			&C.cm[0], &C.sm[0], &C.rmm[0],
			&C.xgrid[0], &C.xgrid[C.kn-1],
			&C.rrr[0], &C.ar[0], &C.ac, &knpar, &idbg)

		for i := 0; i < int(C.kn); i++ {
			fmt.Printf("%3d%8.3f%8.3f\n", i, float64(C.rrr[i]), float64(C.ar[i]))
			C.sd[i] = C.ar[i]
		}
	}
	log.Println("Successful!")

}

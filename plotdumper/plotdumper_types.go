package plotdumper

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/vdobler/chart"
	"github.com/vdobler/chart/imgg"
	"github.com/vdobler/chart/svgg"
	"github.com/vdobler/chart/txtg"

	svg "github.com/ajstarks/svgo"
)

// PlotDumper - структура, описывающая объект, сохраняющий данные
type PlotDumper struct {
	N, M, W, H, Cnt           int
	S                         *svg.SVG
	I                         *image.RGBA
	svgFile, imgFile, txtFile *os.File
}

// NewPlotDumper - конструктор
// n, m - количество графиков по-вертикали и по-горизонтали
// w, h - ширина и высота одной панели
func NewPlotDumper(name string, n, m, w, h int) (*PlotDumper, error) {
	var err error
	dumper := PlotDumper{N: n, M: m, W: w, H: h}

	dumper.svgFile, err = os.Create(name + ".svg")
	if err != nil {
		return nil, fmt.Errorf("error creating svg file %s", err)
	}
	dumper.S = svg.New(dumper.svgFile)
	dumper.S.Start(n*w, m*h)
	dumper.S.Title(name)
	dumper.S.Rect(0, 0, n*w, m*h, "fill: #ffffff")

	dumper.imgFile, err = os.Create(name + ".png")
	if err != nil {
		return nil, fmt.Errorf("error creating png file %s", err)
	}
	dumper.I = image.NewRGBA(image.Rect(0, 0, n*w, m*h))
	bg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 0xff})
	draw.Draw(dumper.I, dumper.I.Bounds(), bg, image.Point{}, draw.Src)

	dumper.txtFile, err = os.Create(name + ".txt")
	if err != nil {
		return nil, fmt.Errorf("error creating txt file %s", err)
	}

	return &dumper, nil
}

// Close - закрывает дампер
func (d *PlotDumper) Close() error {
	err := png.Encode(d.imgFile, d.I)
	if err != nil {
		return fmt.Errorf("проблемы с сохранением файла в png, %s", err)

	}
	err = d.imgFile.Close()
	if err != nil {
		return fmt.Errorf("проблемы с закрытием файла в png, %s", err)

	}
	d.S.End()

	err = d.svgFile.Close()
	if err != nil {
		return fmt.Errorf("проблемы с закрытием файла в svg, %s", err)

	}

	err = d.txtFile.Close()
	if err != nil {
		return fmt.Errorf("проблемы с закрытием файла в txt, %s", err)

	}
	return nil
}

// Plot - отрисовка графика
func (d *PlotDumper) Plot(c chart.Chart) error {
	row, col := d.Cnt/d.N, d.Cnt%d.N

	igr := imgg.AddTo(d.I, col*d.W, row*d.H, d.W, d.H, color.RGBA{0xff, 0xff, 0xff, 0xff}, nil, nil)
	c.Plot(igr)

	sgr := svgg.AddTo(d.S, col*d.W, row*d.H, d.W, d.H, "FiraCode Nerd Font", 12, color.RGBA{0xff, 0xff, 0xff, 0xff})
	c.Plot(sgr)

	tgr := txtg.New(80, 30)
	c.Plot(tgr)
	_, err := d.txtFile.Write([]byte(tgr.String() + "\n\n\n"))
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	d.Cnt++
	return nil
}

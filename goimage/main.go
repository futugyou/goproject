package main

import (
	"math"

	"github.com/fogleman/gg"
)

func main() {
	dc := gg.NewContext(1000, 1000)
	dc.DrawCircle(500, 500, 400)
	//rgb := []float64{255, 222, 173}
	//newRgb := rbgConverter(rgb)
	//dc.SetRGB(1, 1, 0)
	//dc.SetRGB(newRgb[0], newRgb[1], newRgb[2])
	dc.SetRGB255(255, 222, 173)
	dc.Fill()
	dc.SavePNG("demo.jpg")

	im, err := gg.LoadImage("demo.jpg")
	if err != nil {
		panic(err)
	}
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y

	dc = gg.NewContext(h, w)
	redius := math.Min(float64(w), float64(h)) / 2
	dc.DrawRegularPolygon(3, float64(w/2), float64(h/2), redius, redius)

	dc.Clip()
	dc.RotateAbout(gg.Radians(10), float64(w/2), float64(h/2))

	dc.Clear()
	dc.DrawImage(im, 0, 0)

	dc.SetRGB(0, 0, 0)
	s := "Hello, world!"
	sWidth, sHeight := dc.MeasureString(s)
	dc.DrawString(s, (1000-sWidth)/2, (1000+sHeight)/2)

	dc.SavePNG("demo.jpg")
}

func rbgConverter(r []float64) []float64 {
	result := make([]float64, 0)
	for i := 0; i < len(r); i++ {
		result = append(result, r[i]/255)
	}
	return result
}

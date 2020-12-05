package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	//"github.com/Arafatk/glot"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("need data file")
		return
	}
	file := os.Args[1]
	_, err := os.Stat(file)
	if err != nil {
		fmt.Println("cannot stat", file)
		return
	}
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("cannot open", file)
		fmt.Println(err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	allRecords, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	xp := []float64{}
	yp := []float64{}

	for _, rec := range allRecords {
		x, _ := strconv.ParseFloat(rec[0], 64)
		y, _ := strconv.ParseFloat(rec[1], 64)
		xp = append(xp, x)
		yp = append(yp, y)
	}

	points := [][]float64{}
	points = append(points, xp)
	points = append(points, yp)
	fmt.Println(points)

	// dimensions := 2
	// persist := true
	// debug := false
	// plot, _ := glot.NewPlot(dimensions, persist, debug)
	// plot.SetTitle("using glot with csv data")
	// plot.SetXLabel("X-Axis")
	// plot.SetYLabel("Y-Axis")
	// style := "circle"
	// plot.AddPointGroup("Ciecle:", style, points)
	// plot.SavePlot("output.png")
}

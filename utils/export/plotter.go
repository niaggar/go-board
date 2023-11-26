package export

import (
	"bufio"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"os"
	"strconv"
	"strings"
)

func PlotHistogram(saveRoute, dataRoute string) {
	file, err := os.Open(dataRoute)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var data []float64
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		value, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			panic(err)
		}
		data = append(data, value)
	}

	var v plotter.Values
	for i := 0; i < len(data); i++ {
		v = append(v, data[i])
	}

	p := plot.New()
	p.Title.Text = "histogram plot"

	lenght := vg.Length(len(v))

	h, err := plotter.NewBarChart(v, lenght)
	if err != nil {
		panic(err)
	}

	p.Add(h)

	err = p.Save(20*vg.Inch, 30*vg.Inch, saveRoute)
	if err != nil {
		panic(err)
	}
}

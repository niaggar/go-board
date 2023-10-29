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

	p := plot.New()
	v := make(plotter.Values, len(data))
	for i := 0; i < len(data); i++ {
		v[i] = data[i]
	}

	h, err := plotter.NewHist(v, len(data))
	if err != nil {
		panic(err)
	}

	h.Normalize(1)
	p.Add(h)

	err = p.Save(4*vg.Inch, 4*vg.Inch, saveRoute)
	if err != nil {
		panic(err)
	}
}

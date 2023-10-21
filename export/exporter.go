package export

import (
	"fmt"
	"go-board/models"
	"os"
)

type Exporter struct {
	pahtRoute      string
	histogramRoute string
	countElements  int
	totalElements  int
	pathFile       *os.File
	histogramFile  *os.File
	tempPathData   string
}

func NewExporter(pathRoute, histogramRoute string) *Exporter {
	return &Exporter{
		pahtRoute:      pathRoute,
		histogramRoute: histogramRoute,
	}
}

func (e *Exporter) CreatePathFile() {
	file, err := os.Create(e.pahtRoute)
	if err != nil {
		panic(err)
	}
	e.pathFile = file
}

func (e *Exporter) CreateHistogramFile() {
	file, err := os.Create(e.histogramRoute)
	if err != nil {
		panic(err)
	}
	e.histogramFile = file
}

func (e *Exporter) ClosePathFile() {
	err := e.pathFile.Close()
	if err != nil {
		panic(err)
	}
}

func (e *Exporter) CloseHistogramFile() {
	err := e.histogramFile.Close()
	if err != nil {
		panic(err)
	}
}

func (e *Exporter) ExportHistogram() {
	err := e.histogramFile.Sync()
	if err != nil {
		return
	}

	content := fmt.Sprintf("%d \n", e.countElements)
	_, err = e.histogramFile.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func (e *Exporter) ExportElement(sphere models.Sphere) {
	err := e.pathFile.Sync()
	if err != nil {
		return
	}

	content := fmt.Sprintf("%d \t %d \t %f \t %f \t %f \n",
		e.countElements,
		sphere.Type,
		sphere.Position.X,
		sphere.Position.Y,
		sphere.Radius,
	)

	_, err = e.pathFile.WriteString(content)
	if err != nil {
		panic(err)
	}

	e.countElements++
}

func (e *Exporter) StartFrame(total int) {
	e.totalElements = total

	content := fmt.Sprintf("%d\naver\n", e.totalElements)
	_, err := e.pathFile.WriteString(content)
	if err != nil {
		panic(err)
	}

	e.countElements = 0
}

package export

import (
	"fmt"
	"go-board/models"
	"os"
)

type Exporter struct {
	route         string
	countElements int
	totalElements int
	file          *os.File
}

func NewExporter(route string) *Exporter {
	return &Exporter{
		route: route,
	}
}

func (e *Exporter) CreateFile() {
	file, err := os.Create(e.route)
	if err != nil {
		panic(err)
	}
	e.file = file
}

func (e *Exporter) CloseFile() {
	err := e.file.Close()
	if err != nil {
		panic(err)
	}
}

func (e *Exporter) ExportElement(sphere models.Sphere) {
	err := e.file.Sync()
	if err != nil {
		return
	}

	content := fmt.Sprintf("%d \t %d \t %f \t %f \t %f \n",
		e.countElements,
		sphere.Type,
		sphere.Position.X(),
		sphere.Position.Y(),
		sphere.Radius,
	)

	_, err = e.file.WriteString(content)
	if err != nil {
		panic(err)
	}

	e.countElements++
}

func (e *Exporter) StartFrame(total int) {
	e.totalElements = total

	content := fmt.Sprintf("%d\naver\n", e.totalElements)
	_, err := e.file.WriteString(content)
	if err != nil {
		panic(err)
	}

	e.countElements = 0
}

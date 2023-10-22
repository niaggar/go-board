package export

import (
	"bufio"
	"fmt"
	"go-board/gmath"
	"go-board/models"
	"os"
)

type Exporter struct {
	path   string
	file   *os.File
	writer *bufio.Writer
}

func NewExporter(path string) *Exporter {
	return &Exporter{
		path: path,
	}
}

func (e *Exporter) CreateFile() {
	file, err := os.Create(e.path)
	writer := bufio.NewWriterSize(file, 128*1024*4)

	if err != nil {
		panic(err)
	}

	e.file = file
	e.writer = writer
}

func (e *Exporter) CloseFile() {
	e.writer.Flush()
	e.file.Close()
}

func (e *Exporter) Write(content string) {
	_, err := e.writer.WriteString(content)

	if err != nil {
		panic(err)
	}
}

func GetExportPath(number int, sphere *models.Sphere) string {
	content := fmt.Sprintf("%d \t %d \t %f \t %f \t %f \n",
		number,
		sphere.Type,
		sphere.Position.X,
		sphere.Position.Y,
		sphere.Radius,
	)

	return content
}

func GetExportPathBorders(number int, point *gmath.Vector) string {
	content := fmt.Sprintf("%d \t %d \t %f \t %f \t %f \n",
		number,
		2,
		point.X,
		point.Y,
		0.5,
	)

	return content
}

func GetExportHeader(total int) string {
	content := fmt.Sprintf("%d\naver\n", total)
	return content
}

func GetExportHistogram(counts []int) string {
	content := "column\tcount\n"
	for i := 1; i < len(counts); i++ {
		content += fmt.Sprintf("%d\t%d\n", i, counts[i])
	}
	content += "end\n"

	return content
}

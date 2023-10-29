package export

import (
	"fmt"
	"go-board/models"
	"go-board/utils/gmath"
)

type Exporter struct {
	PathFile       *FileWriter
	HistogramFile  *FileWriter
	pathProps      models.ExportProps
	histogramProps models.ExportProps
}

func NewExporter(pathExp, histogramExp models.ExportProps) *Exporter {
	pathFileWriter := NewFileWriter(pathExp.Route)
	histogramFileWriter := NewFileWriter(histogramExp.Route)

	pathFileWriter.CreateFile()
	histogramFileWriter.CreateFile()

	return &Exporter{
		pathProps:      pathExp,
		histogramProps: histogramExp,
		PathFile:       pathFileWriter,
		HistogramFile:  histogramFileWriter,
	}
}

func (e *Exporter) CloseFiles() {
	e.PathFile.CloseFile()
	e.HistogramFile.CloseFile()
}

func (e *Exporter) ExportCurrentState(balls, obstacles []*models.Ball, borders []*gmath.Vector) {
	if !e.pathProps.Active {
		return
	}

	total := len(balls) + len(obstacles) + len(borders)

	content := getExportHeader(total)
	e.PathFile.Write(content)

	itemCount := 0
	for itemCount < len(balls) {
		content = getExportPath(itemCount, balls[itemCount])
		e.PathFile.Write(content)
		itemCount++
	}
	for j := 0; j < len(obstacles); j++ {
		content = getExportPath(itemCount, obstacles[j])
		e.PathFile.Write(content)
		itemCount++
	}
	for j := 0; j < len(borders); j++ {
		content = getExportPathBorders(itemCount, borders[j])
		e.PathFile.Write(content)
		itemCount++
	}
}

func (e *Exporter) ExportHistogram(balls []*models.Ball, columns int, columSize float32) {
	if !e.histogramProps.Active {
		return
	}

	finalCountByColumn := make([]int, columns)

	for i := 0; i < len(balls); i++ {
		pos := int(balls[i].Position.X / columSize)
		finalCountByColumn[pos]++
	}

	content := getExportHistogram(finalCountByColumn)
	e.HistogramFile.Write(content)
}

func getExportPath(number int, sphere *models.Ball) string {
	ballType := 0
	if sphere.Static {
		ballType = 1
	}

	content := fmt.Sprintf("%d \t %d \t %f \t %f \t %f \n",
		number,
		ballType,
		sphere.Position.X,
		sphere.Position.Y,
		sphere.Radius,
	)

	return content
}

func getExportPathBorders(number int, point *gmath.Vector) string {
	content := fmt.Sprintf("%d \t %d \t %f \t %f \t %f \n",
		number,
		2,
		point.X,
		point.Y,
		0.5,
	)

	return content
}

func getExportHeader(total int) string {
	content := fmt.Sprintf("%d\naver\n", total)
	return content
}

func getExportHistogram(counts []int) string {
	content := "column\tcount\n"
	for i := 1; i < len(counts); i++ {
		content += fmt.Sprintf("%d\t%d\n", i, counts[i])
	}

	return content
}

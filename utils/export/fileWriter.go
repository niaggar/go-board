package export

import (
	"bufio"
	"os"
)

type FileWriter struct {
	path   string
	file   *os.File
	writer *bufio.Writer
}

func NewFileWriter(path string) *FileWriter {
	return &FileWriter{
		path: path,
	}
}

func (e *FileWriter) CreateFile() {
	file, err := os.Create(e.path)

	bufferSize := 400 * 1024 * 1024
	writer := bufio.NewWriterSize(file, bufferSize)

	if err != nil {
		panic(err)
	}

	e.file = file
	e.writer = writer
}

func (e *FileWriter) CloseFile() {
	err := e.writer.Flush()
	if err != nil {
		return
	}
	err = e.file.Close()
	if err != nil {
		return
	}
}

func (e *FileWriter) Write(content string) {
	_, err := e.writer.WriteString(content)

	if err != nil {
		panic(err)
	}
}

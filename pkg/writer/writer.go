package writer

import (
	"os"
)

type Writer struct {
	path string
	file *os.File
}

func (writer *Writer) openOrCreateFile() {
	file, err := os.OpenFile(writer.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	writer.file = file

}

func Writer_createReader(path string) *Writer {
	writer := &Writer{path: path}

	writer.openOrCreateFile()

	return writer
}

func (writer *Writer) Writer_writeToFile(content string) {
	_, err := writer.file.WriteString(content)

	if err != nil {
		panic(err)
	}
}

func (writer *Writer) CloseFile() {
	writer.file.Close()
}

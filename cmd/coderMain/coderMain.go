package main

import (
	"l2/pkg/coder"
	"l2/pkg/reader"
	"l2/pkg/writer"
)

func main() {
	reader := reader.Reader_createReader("data/input/test")
	writer := writer.Writer_createReader("data/output/test")

	coder := coder.Coder_createCoder(reader, writer)

	coder.Coder_run()

	writer.CloseFile()
}

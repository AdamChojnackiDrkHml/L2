package main

import (
	"fmt"
	"l2/pkg/coder"
	"l2/pkg/reader"
	"l2/pkg/writer"
	"os"
)

func main() {
	fmt.Println(os.Getwd())
	reader := reader.Reader_createReader("data/input/testy1/test1.bin")
	writer := writer.Writer_createReader("data/output/test")

	coder := coder.Coder_createCoder(reader, writer)

	coder.Coder_run()

	writer.CloseFile()
}

package main

import (
	"fmt"
	"l2/pkg/decoder"
	"l2/pkg/reader"
	"l2/pkg/writer"
	"os"
)

func main() {
	fmt.Println(os.Getwd())
	var readerA *reader.Reader
	var writerA *writer.Writer
	fmt.Println(os.Args)
	if len(os.Args) != 3 {
		readerA = reader.Reader_createReader("../data/output/test")
		writerA = writer.Writer_createReader("../data/outputOutput/test")
	} else {
		readerA = reader.Reader_createReader(os.Args[1])
		writerA = writer.Writer_createReader(os.Args[2])
	}

	decoder := decoder.Decoder_createDecoder(readerA, writerA)

	decoder.Decoder_run()

	writerA.CloseFile()
}

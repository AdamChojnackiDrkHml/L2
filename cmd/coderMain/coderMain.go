package main

import (
	"fmt"
	"l2/pkg/coder"
	"l2/pkg/reader"
	"l2/pkg/writer"
	"os"
)

func main() {
	// fmt.Println(os.Getwd())
	var readerA *reader.Reader
	var writerA *writer.Writer
	// fmt.Println(os.Args)
	if len(os.Args) != 3 {
		readerA = reader.Reader_createReader("../data/input/test")
		writerA = writer.Writer_createReader("../data/output/test")
	} else {
		readerA = reader.Reader_createReader(os.Args[1])
		writerA = writer.Writer_createReader(os.Args[2])
	}

	coder := coder.Coder_createCoder(readerA, writerA)

	coder.Coder_run()

	writerA.CloseFile()

	f1, _ := os.Stat(os.Args[1])
	f2, _ := os.Stat(os.Args[2])
	fmt.Println("Stopień kompresji: ", float64(100.0*(1.0-float64(f2.Size())/float64(f1.Size()))))
	coder.Coder_scanFile(os.Args[1])
	// coder.Coder_avgCodingLenght()
}

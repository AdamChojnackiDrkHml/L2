package main

import (
	"fmt"
	"l2/pkg/reader"
	"l2/pkg/writer"
	"os"
)

func main() {
	fmt.Println(os.Getwd())
	reader := reader.Reader_createReader("data/input/test")
	writer := writer.Writer_createReader("data/output/test")
	patch, readCounter := reader.Reader_readDataPatch()
	writer.Writer_writeToFile(patch)
	fmt.Println(patch)
	fmt.Println(readCounter)
	writer.CloseFile()
}

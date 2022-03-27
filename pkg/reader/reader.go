package reader

import (
	"fmt"
	"os"
)

type Reader struct {
	path               string
	file               *os.File
	PatchSize          int64
	ReadSymbolsCounter int
	IsReading          bool
}

func Print(a string) {
	fmt.Println(a)
}

func (reader *Reader) openFile() {
	file, err := os.Open(reader.path)

	if err != nil {
		panic(err)
	}

	reader.file = file

}

func Reader_createReader(path string) *Reader {
	reader := &Reader{path: path, PatchSize: 64, IsReading: true}

	reader.openFile()

	return reader
}

func (reader *Reader) Reader_readDataPatch() []byte {

	symbols := make([]byte, 0)
	readCounter := 0

	for i := 0; i < 64; i++ {
		currSymbol := make([]byte, 1)
		control, _ := reader.file.Read(currSymbol)

		if control == 0 {
			reader.closeFile()
			reader.IsReading = false
			break
		}

		symbols = append(symbols, currSymbol...)
		readCounter++
	}
	reader.ReadSymbolsCounter = readCounter
	return symbols
}

func (reader *Reader) closeFile() {
	reader.file.Close()
}

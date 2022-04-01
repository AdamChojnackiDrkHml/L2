package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Reader struct {
	path               string
	file               *os.File
	PatchSize          int64
	ReadSymbolsCounter int
	IsReading          bool
	scanner            *bufio.Scanner
	counter            int
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
	reader := &Reader{path: path, PatchSize: 256, IsReading: true}

	reader.openFile()
	reader.counter = 0
	reader.scanner = bufio.NewScanner(reader.file)
	return reader
}

func (reader *Reader) Reader_readDataPatch() []byte {
	currPatch := make([]byte, reader.PatchSize)
	control, err := reader.file.Read(currPatch)

	if err == io.EOF {
		reader.closeFile()
		reader.IsReading = false

	}

	reader.ReadSymbolsCounter = control
	reader.counter += control
	return currPatch[:control]
}

func (reader *Reader) Reader_readLine() []string {
	reader.IsReading = reader.scanner.Scan()

	return strings.Split(reader.scanner.Text(), " ")
}

func (reader *Reader) Reader_readByte() byte {
	oneByteSlice := make([]byte, 1)

	_, err := reader.file.Read(oneByteSlice)

	if err == io.EOF {
		reader.closeFile()
		reader.IsReading = false
		return byte(0)
	}

	return oneByteSlice[0]
}

func (reader *Reader) closeFile() {
	reader.file.Close()
}

func Reader_resetFile(reader **Reader) {
	(*reader) = Reader_createReader((*reader).path)
}

func (reader *Reader) Reader_getFirstWord() string {
	word := make([]byte, 0)
	var char byte

	for reader.IsReading {
		char = reader.Reader_readByte()
		if char == byte(' ') {
			break
		}
		word = append(word, char)
	}
	return string(word)
}

func (reader *Reader) ReadWholeFileGetSizeAndResetReader() int64 {
	var err error
	var counter int64

	for err != io.EOF {
		var control int
		currPatch := make([]byte, reader.PatchSize)
		control, err = reader.file.Read(currPatch)

		reader.ReadSymbolsCounter = control
		counter += int64(control)
	}

	reader.closeFile()
	reader.openFile()

	return int64(counter)
}

package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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
	buffer             []byte
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

func (reader *Reader) Reader_readByte() (byte, error) {
	oneByteSlice := make([]byte, 1)

	_, err := reader.file.Read(oneByteSlice)

	if err == io.EOF {
		reader.closeFile()
		reader.IsReading = false
		return byte(0), err
	}

	return oneByteSlice[0], nil
}

func (reader *Reader) closeFile() {
	reader.file.Close()
}

func Reader_resetFile(reader **Reader) {
	(*reader) = Reader_createReader((*reader).path)
}

// func (reader *Reader) Reader_getFirstWord() string {
// 	word := make([]byte, 0)
// 	var char byte

// 	for reader.IsReading {
// 		char = reader.Reader_readByte()
// 		if char == byte(' ') {
// 			break
// 		}
// 		word = append(word, char)
// 	}
// 	return string(word)
// }

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

func (r *Reader) Reader_ReadBit() byte {

	if len(r.buffer) == 0 {
		r.readByte()
	}

	bit := r.buffer[0]
	r.buffer = r.buffer[1:]

	return bit
}

func (r *Reader) Reader_PeekBit() byte {
	if len(r.buffer) == 0 {
		r.readByte()
	}
	return r.buffer[0]
}

func (r *Reader) Reader_ReadNBits(n int) []byte {
	bits := make([]byte, 0)

	for n != 0 {

		if len(r.buffer) == 0 {
			r.readByte()
		}

		copyBitsUpTo(&r.buffer, &bits, &n)

	}

	return bits
}

func copyBitsUpTo(from *[]byte, to *[]byte, counter *int) {
	howManyCopy := min(len(*from), *counter)

	(*to) = append(*to, (*from)[:howManyCopy]...)
	(*from) = (*from)[howManyCopy:]
	*counter -= howManyCopy
}

func (reader *Reader) readByte() {
	oneByteSlice := make([]byte, 1)

	_, err := reader.file.Read(oneByteSlice)

	if err == io.EOF {
		reader.closeFile()
		reader.IsReading = false
		return
	}
	fmt.Println(oneByteSlice[0])
	reader.buffer = splitByteToBits(oneByteSlice[0])
}

func splitByteToBits(aByte byte) []byte {
	strBits := strings.Split(fmt.Sprintf("%08b", aByte), "")

	bits := make([]byte, 0)

	for _, n := range strBits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		bits = append(bits, byte(bit))
	}
	return bits
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

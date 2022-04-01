package decoder

import (
	"fmt"
	"l2/pkg/reader"
	"l2/pkg/writer"
	"math"
	"strconv"
	"strings"
)

type Decoder struct {
	reader         *reader.Reader
	writer         *writer.Writer
	probsF         []float64
	counterSymbols []int64
	iterations     int
	currentPatch   []byte
	lastPatch      bool
	currTag        []byte
	remOfBit       []byte
	numOfSymbols   int64
	m              int
	bytesBuffer    []byte
	toWrite        string
}

func Decoder_createDecoder(reader *reader.Reader, writer *writer.Writer) *Decoder {
	decoder := &Decoder{reader: reader,
		writer:         writer,
		probsF:         make([]float64, 257),
		counterSymbols: make([]int64, 256),
		iterations:     0,
		currentPatch:   make([]byte, 0),
		lastPatch:      false,
		bytesBuffer:    make([]byte, 0)}

	for i := range decoder.counterSymbols {
		decoder.counterSymbols[i] = 1
	}

	singleProb := 1.0 / 256.0
	for i := range decoder.probsF {
		decoder.probsF[i] = singleProb * float64(i)
	}
	return decoder
}

func (decoder *Decoder) calcProbs() {
	currentPatch := decoder.currentPatch
	decoder.iterations++

	for _, n := range currentPatch {
		decoder.counterSymbols[n]++
	}

	allSymbolsCounter := int64(decoder.iterations+1) * int64(decoder.reader.PatchSize)

	decoder.probsF[0] = 0.0
	for i := 1; i < len(decoder.counterSymbols); i++ {
		temp := float64(decoder.counterSymbols[i-1]) / float64(allSymbolsCounter)
		decoder.probsF[i] = decoder.probsF[i-1] + temp
		//fmt.Println(coder.probsF[i].String())
	}

}

func (decoder *Decoder) decode() {

	decoder.m = int(math.Ceil(math.Log2(1/decoder.probsF[1]))) + 3
	decoder.getMBitsToBuffer()
	fTag := getFloatFromBits(decoder.currTag)
	currChar := byte(0)
	l := 0.0
	p := 1.0

	counter := 0

	for decoder.reader.IsReading && decoder.numOfSymbols > 0 {
		patchSize := 256

		for patchSize != 0 && decoder.numOfSymbols > 0 {

			if p < 0.5 {
				l = l * 2.0
				p = p * 2.0
				decoder.deleteFromBuffer(1 + counter)
				counter = 0
			} else if l >= 0.5 {
				l = (l * 2.0) - 1.0
				p = (p * 2.0) - 1.0
				decoder.deleteFromBuffer(1 + counter)
				counter = 0

			} else if caseThreeCondtionCheck(l, p) {
				l = (l * 2.0) - 0.5
				p = (p * 2.0) - 0.5
				counter++
			} else {
				fTag = getFloatFromBits(decoder.currTag)
				currChar, l, p = decoder.findSymbolWithProbs(fTag, l, p)
				decoder.addByteToBuffer(currChar)
				patchSize--
				decoder.numOfSymbols--
			}
		}

		decoder.calcProbs()
		decoder.calcM()
		decoder.getMBitsToBuffer()
	}

	decoder.toWrite = string(decoder.bytesBuffer)
	decoder.writeCode()
}

func (decoder *Decoder) calcM() {
	minProbs := decoder.probsF[1]
	for i := range decoder.counterSymbols {
		if decoder.probsF[i+1]-decoder.probsF[i] < minProbs {
			minProbs = decoder.probsF[i+1] - decoder.probsF[i]
		}
	}

	newM := int(math.Ceil(math.Log2(1/minProbs))) + 3

	if newM < decoder.m {
		decoder.remOfBit = append(decoder.currTag[decoder.m:], decoder.remOfBit...)
		decoder.currTag = decoder.currTag[:decoder.m]
	} else {
		decoder.getMBitsToBuffer()
	}
}

func (decoder *Decoder) addByteToBuffer(myByte byte) {
	decoder.bytesBuffer = append(decoder.bytesBuffer, myByte)

	if len(decoder.bytesBuffer) == 256 {
		decoder.toWrite = string(decoder.bytesBuffer)
		decoder.writeCode()
		decoder.bytesBuffer = make([]byte, 0)
	}
}
func caseThreeCondtionCheck(l, p float64) bool {
	if p <= 0.5 {
		return false
	}

	if l >= 0.5 {
		return false
	}

	if p >= 0.75 {
		return false
	}

	if l < 0.25 {
		return false
	}

	return true
}

func (decoder *Decoder) writeCode() {
	decoder.writer.Writer_writeToFile(decoder.toWrite)
}

func (decoder *Decoder) getMBitsToBuffer() {

	if len(decoder.currTag) == decoder.m || !decoder.reader.IsReading {
		return
	}

	if len(decoder.remOfBit) != 0 {
		toAdd := decoder.m - len(decoder.currTag)
		toAdd = min(toAdd, len(decoder.remOfBit))
		decoder.currTag = append(decoder.currTag, decoder.remOfBit[:toAdd]...)
		decoder.remOfBit = decoder.remOfBit[toAdd:]

		if len(decoder.currTag) == decoder.m {
			return
		}
	}

	myByte := decoder.reader.Reader_readByte()
	decoder.remOfBit = splitByteToBits(myByte)

	for len(decoder.currTag) != decoder.m && decoder.reader.IsReading {
		if decoder.m-len(decoder.currTag) >= 8 {
			decoder.currTag = append(decoder.currTag, decoder.remOfBit...)
			myByte = decoder.reader.Reader_readByte()
			decoder.remOfBit = splitByteToBits(myByte)
		} else {
			toAdd := decoder.m - len(decoder.currTag)
			decoder.currTag = append(decoder.currTag, decoder.remOfBit[:toAdd]...)
			decoder.remOfBit = decoder.remOfBit[toAdd:]
		}
	}
}

func (decoder *Decoder) findSymbolWithProbs(tag, l, p float64) (byte, float64, float64) {

	d := p - l
	for i := range decoder.counterSymbols {
		high := l + (decoder.probsF[i+1] * d)
		low := l + (decoder.probsF[i] * d)
		if i == 98 {
			fmt.Println(high, " ", low)
		}
		if tag < high && tag >= low {
			return byte(i), low, high
		}
	}
	return 255, decoder.probsF[255], 1.0
}

func (decoder *Decoder) deleteFromBuffer(toDelete int) {
	decoder.currTag = decoder.currTag[toDelete:]
	decoder.getMBitsToBuffer()
}

func (decoder *Decoder) getNumOfSymbols() {

	s := decoder.reader.Reader_getFirstWord()
	if !decoder.reader.IsReading {
		return
	}
	decoder.numOfSymbols, _ = strconv.ParseInt(s, 10, 64)
}

func (decoder *Decoder) Decoder_run() {
	decoder.getNumOfSymbols()
	decoder.decode()
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

func getFloatFromBits(bits []byte) float64 {
	var acc float64

	for i, n := range bits {
		acc += float64(n) * (math.Pow(2, float64(-(i + 1))))
	}
	return acc
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

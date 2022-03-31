package coder

import (
	"encoding/binary"
	"fmt"
	"l2/pkg/reader"
	"l2/pkg/writer"
	"math"
	"os"
	"strconv"
)

type Coder struct {
	reader         *reader.Reader
	writer         *writer.Writer
	probsF         []float64
	counterSymbols []int64
	iterations     int
	currentPatch   []byte
	lastPatch      bool
	tag            float64
	w              string
	bitBuffer      []byte
	bytesBuffer    []byte
}

func Coder_createCoder(reader *reader.Reader, writer *writer.Writer) *Coder {
	coder := &Coder{reader: reader,
		writer:         writer,
		probsF:         make([]float64, 257),
		counterSymbols: make([]int64, 256),
		iterations:     0,
		currentPatch:   make([]byte, 0),
		lastPatch:      false,
		bitBuffer:      make([]byte, 0),
		bytesBuffer:    make([]byte, 0)}

	for i := range coder.counterSymbols {
		coder.counterSymbols[i] = 1
	}

	singleProb := 1.0 / 256.0
	for i := range coder.probsF {
		coder.probsF[i] = singleProb * float64(i)
	}

	return coder
}

func (coder *Coder) calcProbs() {
	currentPatch := coder.currentPatch
	coder.iterations++

	for _, n := range currentPatch {
		coder.counterSymbols[n]++
	}

	allSymbolsCounter := int64(coder.iterations+1) * int64(coder.reader.PatchSize)

	coder.probsF[0] = 0.0
	for i := 1; i < len(coder.counterSymbols); i++ {
		temp := float64(coder.counterSymbols[i-1]) / float64(allSymbolsCounter)
		coder.probsF[i] = coder.probsF[i-1] + temp
		//fmt.Println(coder.probsF[i].String())
	}

}

func (coder *Coder) getData() {
	coder.currentPatch = coder.reader.Reader_readDataPatch()
	coder.lastPatch = !coder.reader.IsReading
}

func (coder *Coder) writeCode() {
	coder.writer.Writer_writeToFile(strconv.Itoa(len(coder.currentPatch)) + " " + fmt.Sprintf("%f", coder.tag) + "\n")
}

func (coder *Coder) writeBytesBuffer() {
	coder.writer.Writer_writeToFile(coder.w)
}

func (coder *Coder) code() {

	l := 0.0
	p := 1.0
	counter := 0
	// iterations := 0
	for !coder.lastPatch {
		coder.getData()
		// fmt.Println(iterations)
		// iterations++
		for _, n := range coder.currentPatch {
			d := p - l
			//fmt.Println(d.String())
			// fmt.Println(coder.probsF)
			// fmt.Println(coder.probsF[int(n)+1])
			p = l + (coder.probsF[int(n)+1] * float64(d))
			//fmt.Println(p.String())

			l = l + (coder.probsF[int(n)] * float64(d))
			//fmt.Println(l.String())
			for {
				if p < 0.5 {
					l = l * 2.0
					p = p * 2.0

					coder.addToBuffer(0)
					for counter > 0 {
						coder.addToBuffer(1)
						counter--
					}

				} else if l >= 0.5 {
					l = (l * 2.0) - 1.0
					p = (p * 2.0) - 1.0
					coder.addToBuffer(1)
					for counter > 0 {
						coder.addToBuffer(0)
						counter--
					}
				} else if caseThreeCondtionCheck(l, p) {
					l = (l * 2.0) - 0.5
					p = (p * 2.0) - 0.5
					counter++
				} else {
					break
				}
			}
		}
		coder.calcProbs()

	}
	temp := p - l
	temp = temp / 2

	coder.tag = l + temp
	tagBits := float64ToByte(coder.tag)

	for _, n := range tagBits {
		coder.addToBuffer(n)

		if n == 1 {
			break
		}
	}

	for len(coder.bitBuffer) != 0 {
		coder.addToBuffer(0)
	}

	if len(coder.bytesBuffer) != 0 {
		coder.w = string(coder.bytesBuffer)
		coder.writeBytesBuffer()
		coder.bytesBuffer = make([]byte, 0)
	}

}

func float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func (coder *Coder) addToBuffer(bit byte) {
	coder.bitBuffer = append(coder.bitBuffer, bit)

	if len(coder.bitBuffer) == 8 {
		coder.addBitsToByteBuffer()
		coder.bitBuffer = make([]byte, 0)
	}
}

func (coder *Coder) addBitsToByteBuffer() {

	acc := byte(0)

	for _, n := range coder.bitBuffer {
		acc *= 2
		acc += n
	}

	coder.addByteToBuffer(acc)
	coder.bitBuffer = make([]byte, 0)

}

func (coder *Coder) addByteToBuffer(myByte byte) {

	coder.bytesBuffer = append(coder.bytesBuffer, myByte)

	if len(coder.bytesBuffer) == 256 {
		coder.w = string(coder.bytesBuffer)
		coder.writeBytesBuffer()
		coder.bytesBuffer = make([]byte, 0)
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

func (coder *Coder) Coder_run() {
	for !coder.lastPatch {
		coder.getData()
		coder.code()
	}
}

func (coder *Coder) Coder_scanFile(path string) {

	counterSlice := make([]int, 256)
	probs1 := make([]float64, 256)

	counter := 0

	f, _ := os.Open(path)

	currSymbol := make([]byte, 1)

	for {
		control, _ := f.Read(currSymbol)
		if control == 0 {
			break
		}
		counter++
		counterSlice[currSymbol[0]]++
	}

	for i, k := range counterSlice {
		probs1[i] = float64(k) / float64(counter)
	}

	H := 0.0

	for i := 0; i < 256; i++ {
		Px := probs1[i]
		if Px != 0.0 {
			Ix := -math.Log2(Px)
			H += Px * Ix
		}
	}

	fmt.Println("ENTROPIA: ", H)
	f.Close()
}

// func (coder *Coder) Coder_avgCodingLenght() {
// 	avg := 0.0
// 	for i, n := range coder.probs {
// 		avg += n * float64(len(coder.codeMap[byte(i)]))
// 	}

// 	fmt.Println("Średnia długość kodowania: ", avg)
// }

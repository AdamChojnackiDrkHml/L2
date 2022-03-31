package coder

import (
	"fmt"
	"l2/pkg/reader"
	"l2/pkg/writer"
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
		coder.w = fmt.Sprintf("%f", coder.probsF[i])
	}

	return coder
}

func (coder *Coder) calcProbs() {
	currentPatch := coder.currentPatch
	iterations := coder.iterations

	for _, n := range currentPatch {
		coder.counterSymbols[n]++
	}

	allSymbolsCounter := int64(iterations)*int64(coder.reader.PatchSize) + int64(coder.reader.ReadSymbolsCounter) + 256

	for i := 0; i < len(coder.counterSymbols); i++ {
		temp := float64(coder.counterSymbols[i]) / float64(allSymbolsCounter)
		coder.probsF[i] = coder.probsF[i-1] + temp
		//fmt.Println(coder.probsF[i].String())
	}

	coder.iterations++
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

	for _, n := range coder.currentPatch {
		d := p - l
		//fmt.Println(d.String())
		//fmt.Println(coder.probsF[n+1])
		p = l + (coder.probsF[n+1] * d)
		//fmt.Println(p.String())

		l = l + (coder.probsF[n] * d)
		//fmt.Println(l.String())
		s := fmt.Sprintf("%f", l)
		s = fmt.Sprintf("%f", coder.probsF[n])
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
				s = fmt.Sprintf("%f", l)
				p = (p * 2.0) - 1.0
				coder.addToBuffer(1)
				for counter > 0 {
					coder.addToBuffer(0)
					counter--
				}
				coder.w = s
			} else if caseThreeCondtionCheck(l, p) {
				l = (l * 2.0) - 0.5
				p = (p * 2.0) - 0.5
				counter++
			} else {
				break
			}
		}
	}
	temp := p - l
	temp = temp / 2

	coder.tag = l + temp

}

func (coder *Coder) addToBuffer(bit byte) {
	coder.bitBuffer = append(coder.bitBuffer, bit)

	if len(coder.bitBuffer) == 256 {
		coder.writeBits()
		coder.bitBuffer = make([]byte, 0)
	}
}

func (coder *Coder) writeBits() {
	for len(coder.bitBuffer) > 0 {
		myByteBits := coder.bitBuffer[:8]
		coder.bitBuffer = coder.bitBuffer[8:]

		acc := byte(0)

		for _, n := range myByteBits {
			acc *= 2
			acc += n
		}

		coder.addByteToBuffer(acc)
	}
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
		coder.writeCode()
		coder.calcProbs()
	}
}

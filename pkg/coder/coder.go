package coder

import (
	"l2/pkg/reader"
	"l2/pkg/writer"

	"github.com/shopspring/decimal"
)

type Coder struct {
	reader         *reader.Reader
	writer         *writer.Writer
	probs          []decimal.Decimal
	counterSymbols []int64
	iterations     int
	currentPatch   []byte
	lastPatch      bool
	toWrite        string
}

func Coder_createCoder(reader *reader.Reader, writer *writer.Writer) *Coder {
	coder := &Coder{reader: reader,
		writer:         writer,
		probs:          make([]decimal.Decimal, 256),
		counterSymbols: make([]int64, 256),
		iterations:     0,
		currentPatch:   make([]byte, 0),
		lastPatch:      false}

	for i := range coder.probs {
		coder.probs[i] = decimal.NewFromInt(1).Div(decimal.NewFromInt(reader.PatchSize))
	}

	return coder
}

func (coder *Coder) calcProbs() {
	currentPatch := coder.currentPatch
	iterations := coder.iterations

	for _, n := range currentPatch {
		coder.counterSymbols[n]++
	}

	allSymbolsCounter := int64(iterations)*int64(coder.reader.PatchSize) + int64(coder.reader.ReadSymbolsCounter)

	for i, n := range coder.counterSymbols {
		coder.probs[i] = decimal.NewFromInt(n).Div(decimal.NewFromInt(allSymbolsCounter))
	}

	coder.iterations++
}

func (coder *Coder) getData() {
	coder.currentPatch = coder.reader.Reader_readDataPatch()
	coder.lastPatch = !coder.reader.IsReading
}

func (coder *Coder) prepareAllDataToCoding() {
	coder.getData()
	coder.calcProbs()
}

func (coder *Coder) writeCode() {
	coder.writer.Writer_writeToFile(coder.toWrite)
}

func (coder *Coder) code() {
	coder.toWrite = string(coder.currentPatch)
}

func (coder *Coder) Coder_run() {
	for !coder.lastPatch {
		coder.prepareAllDataToCoding()
		coder.code()
		coder.writeCode()
	}
}

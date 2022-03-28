package decoder

import (
	"l2/pkg/reader"
	"l2/pkg/writer"
	"strconv"

	"github.com/shopspring/decimal"
)

type Decoder struct {
	reader         *reader.Reader
	writer         *writer.Writer
	probsF         []decimal.Decimal
	counterSymbols []int64
	iterations     int
	currentPatch   []byte
	lastPatch      bool
	tag            decimal.Decimal
	size           int
}

func Decoder_createDecoder(reader *reader.Reader, writer *writer.Writer) *Decoder {
	decoder := &Decoder{reader: reader,
		writer:         writer,
		probsF:         make([]decimal.Decimal, 65),
		counterSymbols: make([]int64, 64),
		iterations:     0,
		currentPatch:   make([]byte, 0),
		lastPatch:      false}

	for i := range decoder.probsF {
		decoder.probsF[i] = decimal.NewFromInt(1).Div(decimal.NewFromInt(reader.PatchSize)).Mul(decimal.NewFromInt(int64(i + 1)))
	}

	decoder.probsF[len(decoder.probsF)-1] = decimal.NewFromInt(1)
	return decoder
}

func (decoder *Decoder) calcProbs() {
	currentPatch := decoder.currentPatch
	iterations := decoder.iterations

	for _, n := range currentPatch {
		decoder.counterSymbols[n]++
	}

	allSymbolsCounter := int64(iterations)*int64(decoder.reader.PatchSize) + int64(decoder.reader.ReadSymbolsCounter)

	decoder.probsF[0] = decimal.NewFromInt(decoder.counterSymbols[0]).Div(decimal.NewFromInt(allSymbolsCounter))

	for i := 1; i < len(decoder.counterSymbols); i++ {
		temp := decimal.NewFromInt(decoder.counterSymbols[i]).Div(decimal.NewFromInt(allSymbolsCounter))
		decoder.probsF[i] = decoder.probsF[i-1].Add(temp)
	}

	decoder.iterations++
}

func (decoder *Decoder) getData() {
	data := decoder.reader.Reader_readLine()
	decoder.tag, _ = decimal.NewFromString(data[1])
	decoder.size, _ = strconv.Atoi(data[0])
	decoder.lastPatch = !decoder.reader.IsReading
}

func (decoder *Decoder) decode() {

}

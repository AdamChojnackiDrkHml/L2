package coder

import (
	"l2/pkg/reader"
	"l2/pkg/writer"
	"strconv"

	"github.com/shopspring/decimal"
)

type Coder struct {
	reader         *reader.Reader
	writer         *writer.Writer
	probsF         []decimal.Decimal
	counterSymbols []int64
	iterations     int
	currentPatch   []byte
	lastPatch      bool
	tag            decimal.Decimal
	w              string
}

func Coder_createCoder(reader *reader.Reader, writer *writer.Writer) *Coder {
	coder := &Coder{reader: reader,
		writer:         writer,
		probsF:         make([]decimal.Decimal, 257),
		counterSymbols: make([]int64, 256),
		iterations:     0,
		currentPatch:   make([]byte, 0),
		lastPatch:      false}

	for i := range coder.counterSymbols {
		coder.counterSymbols[i] = 1
	}

	for i := range coder.probsF {
		coder.probsF[i] = decimal.NewFromInt(1).Div(decimal.NewFromInt(256)).Mul(decimal.NewFromInt(int64(i)))
		coder.w = coder.probsF[i].String()
	}

	coder.probsF[len(coder.probsF)-1] = decimal.NewFromInt(1)
	return coder
}

func (coder *Coder) calcProbs() {
	currentPatch := coder.currentPatch
	iterations := coder.iterations

	for _, n := range currentPatch {
		coder.counterSymbols[n]++
	}

	allSymbolsCounter := int64(iterations)*int64(coder.reader.PatchSize) + int64(coder.reader.ReadSymbolsCounter) + 256

	coder.probsF[0] = decimal.NewFromInt(coder.counterSymbols[0]).Div(decimal.NewFromInt(allSymbolsCounter))

	for i := 1; i < len(coder.counterSymbols); i++ {
		temp := decimal.NewFromInt(coder.counterSymbols[i]).Div(decimal.NewFromInt(allSymbolsCounter))
		coder.probsF[i] = coder.probsF[i-1].Add(temp)
		//fmt.Println(coder.probsF[i].String())
	}

	coder.iterations++
}

func (coder *Coder) getData() {
	coder.currentPatch = coder.reader.Reader_readDataPatch()
	coder.lastPatch = !coder.reader.IsReading
}

func (coder *Coder) writeCode() {
	coder.writer.Writer_writeToFile(strconv.Itoa(len(coder.currentPatch)) + " " + coder.tag.String() + "\n")
}

func (coder *Coder) code() {

	l := decimal.NewFromInt(0)
	p := decimal.NewFromInt(1)
	prefix := make([]int64, 0)
	counter := 0

	for _, n := range coder.currentPatch {
		d := p.Sub(l)
		//fmt.Println(d.String())
		//fmt.Println(coder.probsF[n+1])
		p = l.Add(coder.probsF[n+1].Mul(d))
		//fmt.Println(p.String())

		l = l.Add(coder.probsF[n].Mul(d))
		//fmt.Println(l.String())
		s := l.String()
		s = coder.probsF[n].String()
		for p.LessThan(decimal.NewFromFloat(0.5)) {
			l = l.Mul(decimal.NewFromInt(2))
			p = p.Mul(decimal.NewFromInt(2))
			//fmt.Println(p.String())
			//fmt.Println(l.String())

			for counter > 0 {
				prefix = append(prefix, 1)
				counter--
			}
			prefix = append(prefix, 0)
		}
		s = l.String()
		for decimal.NewFromFloat(0.5).LessThanOrEqual(l) {
			l = l.Mul(decimal.NewFromInt(2)).Sub(decimal.NewFromInt(1))
			s = l.String()
			p = p.Mul(decimal.NewFromInt(2)).Sub(decimal.NewFromInt(1))

			for counter > 0 {
				prefix = append(prefix, 0)
				counter--
			}
			prefix = append(prefix, 1)
			coder.w = s
		}

		for caseThreeCondtionCheck(l, p) {
			l = l.Mul(decimal.NewFromInt(2)).Sub(decimal.NewFromFloat(0.5))
			p = p.Mul(decimal.NewFromInt(2)).Sub(decimal.NewFromFloat(0.5))
			counter++
		}
	}

	coder.tag = l.Add(p)
	coder.tag = coder.tag.Mul(decimal.NewFromFloat(0.5))

	coder.scaleTag(prefix)
}

func (coder *Coder) scaleTag(bits []int64) {

	for _, n := range bits {
		coder.tag = coder.tag.Add(decimal.NewFromInt(1 * n))
		coder.tag = coder.tag.Mul(decimal.NewFromFloat(0.5))

	}
}

func caseThreeCondtionCheck(l, p decimal.Decimal) bool {
	if p.LessThanOrEqual(decimal.NewFromFloat(0.5)) {
		return false
	}

	if decimal.NewFromFloat(0.5).LessThanOrEqual(l) {
		return false
	}

	if decimal.NewFromFloat(0.75).LessThanOrEqual(p) {
		return false
	}

	if l.LessThan(decimal.NewFromFloat(0.25)) {
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

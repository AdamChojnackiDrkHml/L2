package coder

import (
	"fmt"
	"l2/pkg/dictionary"
	"l2/pkg/reader"
	unicode "l2/pkg/uniCode"
	"l2/pkg/writer"
	"math"
	"os"
)

type Coder struct {
	reader  *reader.Reader
	writer  *writer.Writer
	dict    *dictionary.Dictionary
	uniCode unicode.UniCode
}

func Coder_createDefaultCoder(reader *reader.Reader, writer *writer.Writer) *Coder {
	return Coder_createCoder(reader, writer, unicode.Omega)
}

func Coder_createCoder(reader *reader.Reader, writer *writer.Writer, coding unicode.Coding) *Coder {
	coder := &Coder{
		reader: reader,
		writer: writer,
		dict:   dictionary.Dictionary_CreateDictionary()}

	switch coding {
	case unicode.Gamma:
		{
			coder.uniCode = unicode.EliasGamma_Create(reader)
		}
	case unicode.Delta:
		{
			coder.uniCode = unicode.EliasDelta_Create(reader)
		}
	case unicode.Omega:
		{
			coder.uniCode = unicode.EliasOmega_Create(reader)
		}
	case unicode.Fib:
		{
			coder.uniCode = unicode.Fibonacci_Create(reader)
		}
	}

	return coder
}

func (coder *Coder) code() {

	first, err := coder.reader.Reader_readByte()

	if err != nil {
		return
	}

	c := make([]byte, 0)
	c = append(c, first)

	for {
		s, err := coder.reader.Reader_readByte()

		if err != nil {
			break
		}

		if coder.dict.Dictionary_IsContained(append(c, s)) {
			c = append(c, s)
			continue
		}

		codeInt := coder.dict.Dictionary_GetVal(c)
		codeBits := coder.uniCode.CodeNumber(codeInt)
		coder.writer.Writer_addBits(codeBits)

		coder.dict.Dictionary_AddKey(append(c, s))

		c = []byte{s}

	}
	codeInt := coder.dict.Dictionary_GetVal(c)
	codeBits := coder.uniCode.CodeNumber(codeInt)
	coder.writer.Writer_addBits(codeBits)

	coder.writer.Writer_Flush()
}

// func (coder *Coder) writeSize() {
// 	size := coder.reader.ReadWholeFileGetSizeAndResetReader()

// 	coder.w = strconv.FormatInt(size, 10) + " "
// 	coder.writeCode()
// 	coder.w = ""
// }

func (coder *Coder) Coder_run() {
	// coder.writeSize()
	coder.code()

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

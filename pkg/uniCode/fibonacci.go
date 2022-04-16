package unicode

import "l2/pkg/reader"

type Fibonacci struct {
	reader *reader.Reader
	fib    []int
}

func Fibonacci_Create(reader *reader.Reader) *Fibonacci {
	return &Fibonacci{reader: reader, fib: []int{1, 1}}
}

func (f *Fibonacci) getNFib(n int) int {
	last := len(f.fib) - 1
	alLat := last - 1
	for n > len(f.fib)-1 {
		f.fib = append(f.fib, f.fib[last]+f.fib[alLat])
		last++
		alLat++
	}

	return f.fib[n]
}

func (f *Fibonacci) findFibRepresentation(n int) []byte {
	biggestIndex := 0
	bits := make([]byte, 0)
	for n >= int(f.getNFib(biggestIndex)) {
		biggestIndex++
	}
	biggestIndex--
	for ; biggestIndex >= 0; biggestIndex-- {
		if n >= f.fib[biggestIndex] {
			n -= f.fib[biggestIndex]
			bits = append([]byte{1}, bits...)
			continue
		}
		bits = append([]byte{0}, bits...)

	}

	return bits
}

func (f *Fibonacci) getNumberFromFibRepresentation(bits []byte) int {
	acc := 0
	f.getNFib(len(bits) - 1)
	for i, n := range bits {
		acc += f.fib[i] * int(n)
	}

	return acc
}
func (f *Fibonacci) CodeNumber(n int) []byte {
	fibRep := f.findFibRepresentation(n)
	fibRep = append(fibRep, 1)
	return fibRep
}

func (f *Fibonacci) DecodeNumber() int {
	bytes := make([]byte, 0)
	lastBit := byte(0)
	currBit := f.reader.Reader_ReadBit()
	for !(lastBit == 1 && currBit == 1) {
		bytes = append(bytes, currBit)
		lastBit = currBit
		currBit = f.reader.Reader_ReadBit()
	}
	return f.getNumberFromFibRepresentation(bytes)
}

package unicode

import (
	"l2/pkg/reader"
)

type EliasDelta struct {
	reader *reader.Reader
}

func EliasDelta_Create(reader *reader.Reader) *EliasDelta {
	return &EliasDelta{reader: reader}
}

func (eD *EliasDelta) CodeNumber(num int) []byte {
	binary := numberToBits(num)
	prefixN := numberToBits(len(binary))
	binary = binary[1:]
	prefix := make([]byte, len(prefixN)-1)

	prefixN = append(prefixN, binary...)
	prefix = append(prefix, prefixN...)

	return prefix
}

func (ed *EliasDelta) DecodeNumber() int {
	firstCounter := ed.countK()
	readNBits := ed.getBitsOfNumber(firstCounter - 1)

	N := getByteFromBits(readNBits)

	readBits := ed.getBitsOfNumber(N - 1)
	return getByteFromBits(readBits)
}

func (ed *EliasDelta) countK() int {
	counter := 0

	for {
		if ed.reader.Reader_ReadBit() == byte(1) {
			break
		}
		counter++
	}

	return counter
}

func (ed *EliasDelta) getBitsOfNumber(counter int) []byte {
	readBits := make([]byte, 1)
	readBits[0] = 1
	//Copy representation from reminder up to $counter bits
	readBits = append(readBits, ed.reader.Reader_ReadNBits(counter)...)
	return readBits
}

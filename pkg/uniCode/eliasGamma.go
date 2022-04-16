package unicode

import (
	"l2/pkg/reader"
)

type EliasGamma struct {
	reader *reader.Reader
}

func EliasGamma_Create(reader *reader.Reader) *EliasGamma {
	return &EliasGamma{reader: reader}
}

func (eG *EliasGamma) CodeNumber(num int) []byte {
	binary := numberToBits(num)
	prefix := make([]byte, len(binary)-1)

	prefix = append(prefix, binary...)

	return prefix
}

func (eg *EliasGamma) DecodeNumber() int {
	counter := eg.countN()
	readBits := eg.getBitsOfNumber(counter)

	return getByteFromBits(readBits)
}

func (eg *EliasGamma) getBitsOfNumber(counter int) []byte {
	readBits := make([]byte, 1)
	readBits[0] = 1
	//Copy representation from reminder up to $counter bits
	readBits = append(readBits, eg.reader.Reader_ReadNBits(counter)...)
	return readBits
}

func (eg *EliasGamma) countN() int {
	counter := 0

	for {
		if eg.reader.Reader_ReadBit() == byte(1) {
			break
		}
		counter++
	}

	return counter
}

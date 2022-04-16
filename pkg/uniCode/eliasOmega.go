package unicode

import (
	"l2/pkg/reader"
)

type EliasOmega struct {
	reader *reader.Reader
}

func EliasOmega_Create(reader *reader.Reader) *EliasOmega {
	return &EliasOmega{reader: reader}
}

func (eO *EliasOmega) CodeNumber(num int) []byte {
	buffer := make([]byte, 0)
	buffer = append(buffer, 0)
	K := num
	for K != 1 {
		kBits := numberToBits(K)
		buffer = append(kBits, buffer...)

		K = len(kBits) - 1
	}

	return buffer
}

func (eO *EliasOmega) DecodeNumber() int {

	n := 1

	for eO.reader.Reader_PeekBit() != 0 {
		n = getByteFromBits(eO.reader.Reader_ReadNBits(n + 1))
	}
	eO.reader.Reader_ReadBit()
	return n
}

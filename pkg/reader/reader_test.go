package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader_ReadBit(t *testing.T) {

	r := Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/reader/testReader")

	assert.Equal(t, byte(0), r.Reader_ReadBit(), "Not working")

}

func TestReader_ReadNBits(t *testing.T) {

	r := Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/reader/testReader")

	assert.Equal(t, []byte{0, 0, 1, 1, 1, 1, 1, 1}, r.Reader_ReadNBits(8), "Not working")

	r = Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/reader/testReader2")

	assert.Equal(t, []byte{0, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1}, r.Reader_ReadNBits(11), "Not working")

}

func TestReader_TestCombinedRead(t *testing.T) {

	r := Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/reader/testReader")

	assert.Equal(t, byte(0), r.Reader_ReadBit(), "Not working 1")
	assert.Equal(t, []byte{0, 1, 1, 1, 1, 1, 1}, r.Reader_ReadNBits(7), "Not working 2")
}

package unicode

import (
	"l2/pkg/reader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGammaCode(t *testing.T) {
	eG := &EliasGamma{}

	converted := eG.CodeNumber(4)
	converted2 := eG.CodeNumber(127)

	assert.Equal(t, []byte{0, 0, 1, 0, 0}, converted, "Not equal 1")
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1}, converted2, "Not equal 2")

}

func TestGammaDecode(t *testing.T) {
	myReader := reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Gamma/test1byte")

	eG := EliasGamma_Create(myReader)

	assert.Equal(t, 15, eG.DecodeNumber(), "Not equal 1")

	myReader = reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Gamma/test2byte")

	eG = EliasGamma_Create(myReader)

	assert.Equal(t, 15, eG.DecodeNumber(), "Not equal 2")
	assert.Equal(t, 31, eG.DecodeNumber(), "Not equal 3")

}

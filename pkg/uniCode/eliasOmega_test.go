package unicode

import (
	"l2/pkg/reader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOmegaCode(t *testing.T) {
	eO := &EliasOmega{}

	converted := eO.CodeNumber(15)
	converted2 := eO.CodeNumber(127)

	assert.Equal(t, []byte{1, 1, 1, 1, 1, 1, 0}, converted, "Not equal 1")
	assert.Equal(t, []byte{1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 0}, converted2, "Not equal 2")

}

func TestOmegaDecode(t *testing.T) {
	myReader := reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Omega/test1byte")

	eO := EliasOmega_Create(myReader)

	assert.Equal(t, 5, eO.DecodeNumber(), "Not equal 1")

	myReader = reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Omega/test2byte")

	eO = EliasOmega_Create(myReader)

	assert.Equal(t, 15, eO.DecodeNumber(), "Not equal 2")

	myReader = reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Omega/testMany")

	eO = EliasOmega_Create(myReader)

	assert.Equal(t, 5, eO.DecodeNumber(), "Not equal 3")
	assert.Equal(t, 15, eO.DecodeNumber(), "Not equal 4")

}

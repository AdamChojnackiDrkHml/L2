package unicode

import (
	"l2/pkg/reader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeltaCode(t *testing.T) {
	eD := &EliasDelta{}

	converted := eD.CodeNumber(4)
	converted2 := eD.CodeNumber(127)

	assert.Equal(t, []byte{0, 1, 1, 0, 0}, converted, "Not equal 1")
	assert.Equal(t, []byte{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1}, converted2, "Not equal 2")

}

func TestDeltaDecode(t *testing.T) {
	myReader := reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Delta/test1byte")

	eD := EliasDelta_Create(myReader)

	assert.Equal(t, 5, eD.DecodeNumber(), "Not equal 1")

	myReader = reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Delta/test2byte")

	eD = EliasDelta_Create(myReader)

	assert.Equal(t, 15, eD.DecodeNumber(), "Not equal 2")

	myReader = reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Delta/testMany")

	eD = EliasDelta_Create(myReader)

	assert.Equal(t, 5, eD.DecodeNumber(), "Not equal 3")
	assert.Equal(t, 15, eD.DecodeNumber(), "Not equal 4")

}

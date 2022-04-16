package coder

import (
	"l2/pkg/reader"
	unicode "l2/pkg/uniCode"
	"l2/pkg/writer"
	"testing"
)

func TestCode(t *testing.T) {
	r := reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/coder/testFiles/input1")
	w := writer.Writer_createWriter("/Users/Adam/Desktop/KKD/L3/pkg/coder/testFiles/output1")

	c := Coder_createCoder(r, w, unicode.Gamma)

	c.Coder_run()

}

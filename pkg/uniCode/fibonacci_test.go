package unicode

import (
	"l2/pkg/reader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFibonacciCode(t *testing.T) {
	fib := Fibonacci_Create(nil)

	converted := fib.CodeNumber(4)
	converted2 := fib.CodeNumber(127)
	// 1 1 2 3 5 8 13 21 34 55 89 144
	// 0 1 0 1 0 0 0  0  1  0  1  0
	assert.Equal(t, []byte{0, 1, 0, 1, 1}, converted, "Not equal 1")
	assert.Equal(t, []byte{0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1}, converted2, "Not equal 2")

}

func TestFibonacciDecode(t *testing.T) {
	myReader := reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Fibonacci/test1byte")

	fib := Fibonacci_Create(myReader)

	assert.Equal(t, 5, fib.DecodeNumber(), "Not equal 1")

	myReader = reader.Reader_createReader("/Users/Adam/Desktop/KKD/L3/pkg/uniCode/testFiles/Fibonacci/test2byte")

	fib = Fibonacci_Create(myReader)

	assert.Equal(t, 5, fib.DecodeNumber(), "Not equal 3")
	assert.Equal(t, 15, fib.DecodeNumber(), "Not equal 4")

}

package unicode

import (
	"fmt"
	"strconv"
	"strings"
)

type Coding int

const (
	Gamma Coding = iota
	Delta
	Omega
	Fib
)

type UniCode interface {
	CodeNumber(num int) []byte
	DecodeNumber() int
}

func Code(num int, uC UniCode) []byte {
	return uC.CodeNumber(num)
}

func Decode(uC UniCode) int {
	return uC.DecodeNumber()
}

func numberToBits(num int) []byte {
	strBits := strings.Split(fmt.Sprintf("%b", num), "")

	bits := make([]byte, 0)

	for _, n := range strBits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		bits = append(bits, byte(bit))
	}
	return bits
}

func getByteFromBits(bits []byte) int {
	acc := 0

	for _, n := range bits {
		acc *= 2
		acc += int(n)
	}

	return acc
}

func decodeBin(binary []byte) int {
	acc := 0

	for _, n := range binary {
		acc *= 2
		acc += int(n)
	}

	return acc
}

package main

import (
	"encoding/hex"
	"os"
	"unsafe"
)

func main() {
	a := "<--HEXDSC-->"
	b, err := hex.DecodeString(a)
	if err != nil {
		os.Exit(1)
	}

	bb(aa(b))
}

func aa(a []byte) []byte {
	b := a
	return b
}

func bb(a []byte) {
	aa := *(*func() int)(unsafe.Pointer(&a[0]))
	aa()
}

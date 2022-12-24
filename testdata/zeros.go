
// +build ignore

//go:generate go run zeros.go

// Generates zeros.bin. This test file contains zeroed data throughout and
// tests the best case compression scenario.
package main

import "io/ioutil"

const (
	name = "zeros.bin"
	size = 1 << 18
)

func main() {
	b := make([]byte, size)
	if err := ioutil.WriteFile(name, b[:size], 0664); err != nil {
		panic(err)
	}
}

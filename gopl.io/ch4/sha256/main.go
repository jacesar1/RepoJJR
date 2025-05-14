// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 83.

// The sha256 command computes the SHA256 hash (an array) of a string.
package main

import (
	"crypto/sha256" //!+
	"fmt"
)

var pc [256]byte

func init() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

func main() {
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))

	//resutado do XOR
	var bitsDiferentes [sha256.Size]byte
	somatorio := 0

	for i := range sha256.Size {

		bitsDiferentes[i] = c1[i] ^ c2[i]

		somatorio += int(pc[bitsDiferentes[i]])
		x := int(pc[bitsDiferentes[i]])

		fmt.Printf("c1:%08b c2:%08b bitsDiferentes:%08b tem %d bits diferentes\n", c1[i], c2[i], bitsDiferentes[i], x)
	}
	//fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)
	fmt.Printf("%d\n", somatorio)
	fmt.Println('x', 'X')

	// Output:
	// 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881
	// 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
	// false
	// [32]uint8
}

//!-

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 6.
//!+

// Echo2 prints its command-line arguments.
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()
	for i, arg := range os.Args[1:] {
		fmt.Printf("O indice é: %d e o argumento é: %s\n", i, arg)
		//s += sep + arg
		//sep = " "
	}
	fmt.Printf("%.6fs elapsed\n", time.Since(start).Seconds())
}

//!-

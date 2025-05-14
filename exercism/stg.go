package main

import (
	
	"fmt"
)

func Welcome(name string) string {
	return "Welcome to my party, " + name + "!"
}

func main() {
	// Example usage of the Welcome function	
	fmt.Println(Welcome("Alice"))
}


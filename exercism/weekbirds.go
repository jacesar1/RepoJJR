package main

import (
	"fmt"
)


// TotalBirdCount return the total bird count by summing
// the individual day's counts.
func TotalBirdCount(birdsPerDay []int) int {
    soma :=0
	for _, b := range birdsPerDay {
			soma += b
        
    }
    return soma
}

// BirdsInWeek returns the total bird count by summing
// only the items belonging to the given week.
func BirdsInWeek(birdsPerDay []int, week int) int {
	 start := (week - 1) * 7
     end := start + 7
    
     return TotalBirdCount(birdsPerDay[start:end])
 
}

func	main() {
	birdsPerDay := []int{4, 7, 3, 2, 1, 1, 2, 0, 2, 3, 2, 7, 1, 3, 0, 6, 5, 3, 7, 2, 3}
	fmt.Println(TotalBirdCount(birdsPerDay))
	fmt.Println(BirdsInWeek(birdsPerDay, 3))
}
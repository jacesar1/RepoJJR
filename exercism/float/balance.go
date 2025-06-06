package main

import (	

	"fmt"
)

// InterestRate returns the interest rate for the provided balance.
func InterestRate(balance float64) float32 {
   var irate float32
	switch {
        case balance < 0.0 :
        irate = 3.213
        case balance >= 0.0 && balance < 1000.0:
        irate = 0.5
        case balance >= 1000.0 && balance < 5000.0:
        irate = 1.621
        case balance >= 5000.0 :
        irate = 2.475      
       
    }
    return irate
}

// Interest calculates the interest for the provided balance.
func Interest(balance float64) float64 {
   
    return float64(InterestRate(balance)) * balance / 100
}

// AnnualBalanceUpdate calculates the annual balance update, taking into account the interest rate.
func AnnualBalanceUpdate(balance float64) float64 {
	return Interest(balance) + balance
}

// YearsBeforeDesiredBalance calculates the minimum number of years required to reach the desired balance.
func YearsBeforeDesiredBalance(balance, targetBalance float64) int {

    resultado := 0.0
    anos := 0
    for resultado <= targetBalance {
     resultado += Interest(balance) + balance
	 balance = resultado
	 fmt.Printf("Balance: %f\n", resultado)

        anos++
    }
      return anos 
}

func main() {
	// Example usage
	initialBalance := 1000.0
	targetBalance := 2000.0

	years := YearsBeforeDesiredBalance(initialBalance, targetBalance)
	fmt.Printf("Years required to reach %.2f from %.2f: %d\n", targetBalance, initialBalance, years)
}			
package main

import (
	"fmt"
)

// PreparationTime calculates the preparation time for a lasagna based on the number of layers
func PreparationTime(layers []string, min int) int {
	if min == 0 {
		return len(layers) * 2
	}
	return len(layers) * min
}

func Quantities(layers []string) (int, float64) {
     qnoodles := 0
     qsauce := 0.0
    
     for _, layer := range layers{
         if layer == "noodles" {
             qnoodles += 50
         }
         if layer == "sauce" {
             qsauce += 0.2
         }
     }

    return qnoodles, qsauce
}

func AddSecretIngredient(layers, layers2 []string) {

    layers2[len(layers2)-1] = layers[len(layers)-1]
}

func ScaleRecipe(quantities []float64, portions int) []float64 {

    ScaledQuantities := []float64{}
    for i, _ := range quantities {

        ScaledQuantities[i] = quantities[i] * float64((portions/2))
    }
    return ScaledQuantities
}

func main() {
	layers := []string{"noodles", "sauce", "cheese"}
	prepTime := PreparationTime(layers, 0)
	fmt.Println("Preparation time:", prepTime)

	noodles, sauce := Quantities(layers)
	fmt.Println("Noodles (grams):", noodles)
	fmt.Println("Sauce (liters):", sauce)
	
	friend := []string{"noodles", "sauce", "cheese", "nutmeg"}
    mine := []string{"noodles", "sauce", "cheese", "?"}
    AddSecretIngredient(friend, mine)
	fmt.Println("My lasagna layers:", mine)
	
}

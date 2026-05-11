
package main

import (
	"fmt"
)


func addManyNums(nums ...float64) float64{
	store := 0.0

	for _, digit := range nums {
		store += digit
	}
	return store
	
}
func main() {
	store := addManyNums(10,30,30,50,60,30.49,28.48)
	fmt.Println(store)
}

/*
	Declaring variables and concatenating them.
	There are different ways of declaring variables in golang.
*/

package main

import (
	"fmt"
)

func main() {
	var age int = 28

	var (
		height        = 1.78
		name   string = "Luiz Rafael"
	)

	from := "Bahia"
	history := false
	father, mother, ageFather := "Luiz", "Cida", 58

	fmt.Println("My name is: "+name+" I have ", age, "old.")
	fmt.Println("I am from", from, ", my height is ", height, "m")
	fmt.Println(father, mother, "are my parents, the age of my father is ", ageFather)
	fmt.Println("This history is ", history)
}

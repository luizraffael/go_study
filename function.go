package main

import "fmt"

func sum(a int, b int) int {
	return a + b
}

func avg(c, d float32) float32 {
	return (c + d) / 2
}

//You can return two values in goLang
func calc(a, b float32) (float32, float32) {
	var avg = (a + b) / 2
	var sum = a + b
	return avg, sum
}

func nackedReturn(a int) (square int, cube int) {
	square = a * a
	cube = a * a * a
	return
}

func main() {
	a, s := calc(10, 8)

	sqr, c := nackedReturn(2)
	fmt.Println("The result of sum is: ", sum(3, 6))
	fmt.Println("The result of avg is: ", avg(6.5, 9))
	fmt.Println("Returning two values in one function", a, " - ", s)
	fmt.Println("Function nacked return, values: ", sqr, "and", c)
}

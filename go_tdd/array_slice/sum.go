package array_slice

func Sum(numbers []int) int {
	sum := 0
	for _, numbers := range numbers {
		sum += numbers
	}
	return sum
}

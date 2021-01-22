package array_slice

import "testing"

func TestSoma(t *testing.T) {
	numbers := [5]int{1, 2, 3, 4, 5}

	result := Sum(numbers)
	expected := 15

	if expected != result {
		t.Errorf("Result %d, expected %d, data %v", result, expected, numbers)
	}
}

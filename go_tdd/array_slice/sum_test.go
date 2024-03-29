package array_slice

import (
	"reflect"
	"testing"
)

func TestSoma(t *testing.T) {
	t.Run("Colection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		result := Sum(numbers)
		expected := 15

		if expected != result {
			t.Errorf("Result %d, expected %d, data %v", result, expected, numbers)
		}
	})

	t.Run("Dinamic collection", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		result := Sum(numbers)
		expected := 6

		if expected != result {
			t.Errorf("Result %d, expected %d, data %v", result, expected, numbers)
		}
	})
}
func TestSumAll(t *testing.T) {

	received := SumAll([]int{1, 2}, []int{0, 9})
	expectec := []int{3, 9}

	if !reflect.DeepEqual(received, expectec) {
		t.Errorf("result %v expectec %v", received, expectec)
	}
}

package structs_methods_interface

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	result := Perimeter(rectangle)
	expected := 40.0

	if result != expected {
		t.Errorf("Result %.2f expected %.2f", result, expected)
	}
}

func TestArea(t *testing.T) {
	t.Run("rectangles", func(t *testing.T) {
		rectangle := Rectangle{10.0, 10.0}
		result := rectangle.Area(rectangle)
		expected := 100.0

		if result != expected {
			t.Errorf("Resul %.2f expexted %.2f", result, expected)
		}
	})

	t.Run("circles", func(t *testing.T) {
		circle := Circle{10.0}
		result := circle.Area(circle)
		expected := 314.1592653589793

		if result != expected {
			t.Errorf("Resul %.2f expexted %.2f", result, expected)
		}
	})

}

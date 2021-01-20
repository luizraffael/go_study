package main 

import "testing"

func TestHello(t *testing.T){

	verifyCorrectMessage := func(t *testing.T, result, expected string){
		t.Helper()
		if result != expected{
			t.Errorf("result '%s', expected '%s'", result, expected)
		}	
	}

	t.Run("Say hello to people", func(t *testing.T){
		result := Hello("Chris","")
		expected := "Hello, Chris"
		verifyCorrectMessage(t, result, expected)
	
	})

	t.Run("'Hello, world' when 'string' is empty", func(t *testing.T) {
		result := Hello("","")
		expected := "Hello, world"
		verifyCorrectMessage(t, result, expected)

	})

	t.Run("In Spanish", func(t *testing.T){
		result := Hello("Elodie", "Spanish")
		expected := "Hola, Elodie"
		verifyCorrectMessage(t, result, expected)
	})

	t.Run("in French", func(t *testing.T){
		result := Hello("John", "French")
		expected := "Bonjour, John"
		verifyCorrectMessage(t, result, expected)
	})
	
}
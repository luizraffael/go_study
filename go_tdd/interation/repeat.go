package interation

func Repeat(character string) string {
	var repetitions string
	for i := 0; i < 5; i++ {
		repetitions = repetitions + character
	}
	return repetitions
}

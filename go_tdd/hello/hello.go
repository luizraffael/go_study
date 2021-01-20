package main

import "fmt"

const prefixEnlish = "Hello, "
const prefixSpanish = "Hola, "
const prefixFrench = "Bonjour, "
const spanish = "Spanish"
const french = "French"



func Hello(name string, language string) string {
	if name == ""{
		name = "world"
	}

	return prefixGreeting(language) + name
	
}

func prefixGreeting(language string) (prefix string){
	switch language{
	case french: 
		prefix = prefixFrench
	case spanish:
		prefix = prefixSpanish
	default:
		prefix = prefixEnlish
	}
		
	return
}

func main(){
	fmt.Println(Hello("",""))
}
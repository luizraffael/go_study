package main

import "fmt"

const prefixHello = "Hello, "

func Hello(name string) string {
	if name == ""{
		name = "world"
	}
	return prefixHello + name
}

func main(){
	fmt.Println(Hello(""))
}
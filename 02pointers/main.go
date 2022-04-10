package main

import "fmt"

func main() {
	var num int = 21
	fmt.Println("print my number variable", num)

	var count int = 35
	var ptr = &count
	*ptr = *ptr + 2

	fmt.Println("print my pointer", count)
	fmt.Println("print my pointer", *ptr)

}

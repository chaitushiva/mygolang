package main

import (
	"fmt"
	"reflect"
)

func main() {

	var overall [][]int
	fmt.Println("printing overall value and its type %T", overall)
	fmt.Println("overall kind = ", reflect.ValueOf(overall).Kind())
	fmt.Println("overall type = ", reflect.TypeOf(overall))

}

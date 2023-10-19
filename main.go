package main

import (
	"fmt"
	"go-board/bmath"
)

func main() {
	fmt.Println("Hello, world!")

	v := bmath.NewVector(1, 2)
	fmt.Println(v)

	v.SetX(3)
	fmt.Println(v)
}

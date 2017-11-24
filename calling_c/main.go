package main

// #include "c/add.c"
import "C"

import (
	"fmt"
)

func main() {
	r := C.add(40, 2)
	fmt.Println("result = ", r)
}

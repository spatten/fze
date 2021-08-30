package main

import (
	"fmt"
	"fze"
	"os"
)

func main() {
	args := os.Args
	err := fze.Runner(args[1:])
	if err != nil {
		fmt.Printf("Error! %v\n%v\n", args, err)
	}
}

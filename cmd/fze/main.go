package main

import (
	"fmt"
	"fze"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	var runnerOpts fze.RunnerOptions
	if args[0] == "-m" || args[0] == "--multi" {
		args = args[1:]
		runnerOpts.Multi = true
	}
	log.Printf("args: %v", args)

	err := fze.Runner(args, runnerOpts)
	if err != nil {
		fmt.Printf("Error! %v\n%v\n", args, err)
	}
}

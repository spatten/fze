package main

import (
	"flag"
	"fmt"
	"fze"
	"os"
)

func main() {
	args := os.Args[1:]
	var multi bool

	flag.BoolVar(&multi, "multi", false, "Run FZF in multi mode")
	flag.Parse()
	runnerOpts := fze.RunnerOptions{Multi: multi}
	// Everything else is an arg
	args = args[flag.NFlag():]

	err := fze.Runner(args, runnerOpts)
	if err != nil {
		fmt.Printf("Error! %v\n%v\n", args, err)
	}
}

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
	var emacsServer string

	flag.BoolVar(&multi, "multi", false, "Run FZF in multi mode")
	flag.StringVar(&emacsServer, "emacs-server", os.ExpandEnv("$FZE_EMACS_SERVER"), "Emacs Server name (defaults to FZE_EMACS_SERVER env variable)")
	flag.Parse()
	runnerOpts := fze.RunnerOptions{
		Multi:       multi,
		EmacsServer: emacsServer,
	}
	// Everything else is an arg
	args = args[flag.NFlag():]

	err := fze.Runner(args, runnerOpts)
	if err != nil {
		fmt.Printf("Error! %v\n%v\n", args, err)
	}
}

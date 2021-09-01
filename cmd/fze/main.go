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
	var switchToTmuxPane bool
	var emacsServer string

	flag.BoolVar(&multi, "multi", false, "Run FZF in multi mode")
	flag.BoolVar(&switchToTmuxPane, "switch-to-tmux-pane", os.ExpandEnv("$FZE_SWITCH_TO_TMUX_PANE") != "", "Switch to a tmux pane after opening the file")
	flag.StringVar(&emacsServer, "emacs-server", os.ExpandEnv("$FZE_EMACS_SERVER"), "Emacs Server name (defaults to FZE_EMACS_SERVER env variable)")
	flag.Parse()
	runnerOpts := fze.RunnerOptions{
		Multi:            multi,
		EmacsServer:      emacsServer,
		SwitchToTmuxPane: switchToTmuxPane,
	}
	err := fze.Runner(flag.Args(), runnerOpts)
	if err != nil {
		fmt.Printf("Error! %v\n%v\n", args, err)
	}
}

package fze

import (
	"fmt"
	"os/exec"
)

func findRunner(args []string, opts RunnerOptions) error {
	// Get the output from find
	cmd := "find " + fixArgs(args)
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("running find: %v", err)
	}

	// Run the output from find through fzf
	outLines, err := runFzf(res, opts)
	if err != nil {
		return fmt.Errorf("runFzf: %v", err)
	}

	var paths []pathArg
	for _, line := range outLines {
		paths = append(paths, pathArg{path: line})
	}

	// Run emacsclient
	return openEditor(paths, opts)
}

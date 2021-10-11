package fze

import (
	"fmt"
	"os/exec"
)

func fzfRunner(args []string, opts RunnerOptions) error {
	cmd := "find * -type f"
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err
	}

	outLines, err := runFzf(res, opts)
	if err != nil {
		return fmt.Errorf("runFzf: %v", err)
	}

	var paths []pathArg
	for _, line := range outLines {
		paths = append(paths, pathArg{path: line})
	}

	// Run emacsclient
	err = openEditor(paths, opts)
	if err != nil {
		return fmt.Errorf("running emacsclient: %v", err)
	}

	return nil

}

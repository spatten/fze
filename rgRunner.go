package fze

import (
	"fmt"
	"os/exec"
	"strings"
)

func rgRunner(args []string, opts RunnerOptions) error {
	// Get the output from rg
	cmd := "rg -n " + fixArgs(args)
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("running rg: %v", err)
	}

	// Run the output from rg through fzf
	outLines, err := runFzf(res, opts)
	if err != nil {
		return fmt.Errorf("runFzf: %v", err)
	}

	var paths []pathArg
	for _, line := range outLines {
		// Get the filename and linenumber from the output
		output := strings.Split(line, ":")
		if len(output) < 2 {
			return fmt.Errorf("expecting a path and a line-number in this rg output: %s", output)
		}
		paths = append(paths, pathArg{path: output[0], lineNumber: output[1]})
	}

	// Run emacsclient
	return openEditor(paths, opts)
}

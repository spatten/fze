package fze

import (
	"fmt"
	"os/exec"
	"strings"
)

func gitRunner(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("no args to git")
	}

	switch arg := args[0]; arg {
	case "diff":
		return gitDiffRunner(args[1:])
	case "show":
		return gitShowRunner(args[1:])
	}

	return "", fmt.Errorf("only git diff and git show are supported")
}

func gitDiffRunner(args []string) (string, error) {
	// Get the output from git
	cmd := "git diff --src-prefix=a/ --dst-prefix=b/ --color=always " + strings.Join(args, " ") + " | showlinenum.awk show_path=1"
	fmt.Printf("Running cmd: %s\n", cmd)
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("running git: %v", err)
	}

	// Run the output from git through fzf
	out, err := runFzf(res)
	if err != nil {
		return "", fmt.Errorf("runFzf: %v", err)
	}

	// Get the filename and linenumber from the output
	output := strings.Split(out, ":")
	if len(output) < 2 {
		return "", fmt.Errorf("expecting a path and a line-number in this git output: %s", output)
	}
	path := output[0]
	lineNumber := output[1]

	// Run emacsclient
	err = openEditorWithLineNumber(path, lineNumber)
	if err != nil {
		return "", fmt.Errorf("running emacsclient: %v", err)
	}

	return "", nil
}

func gitShowRunner(args []string) (string, error) {
	// Get the output from git
	cmd := "git show --src-prefix=a/ --dst-prefix=b/ --color=always --oneline " + strings.Join(args, " ") + "| tail +2 | showlinenum.awk show_path=1 "
	fmt.Printf("Running cmd: %s\n", cmd)
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("running git: %v", err)
	}

	// Run the output from git through fzf
	out, err := runFzf(res)
	if err != nil {
		return "", fmt.Errorf("runFzf: %v", err)
	}

	// Get the filename and linenumber from the output
	output := strings.Split(out, ":")
	if len(output) < 2 {
		return "", fmt.Errorf("expecting a path and a line-number in this git output: %s", output)
	}
	path := output[0]
	lineNumber := output[1]

	// Run emacsclient
	err = openEditorWithLineNumber(path, lineNumber)
	if err != nil {
		return "", fmt.Errorf("running emacsclient: %v", err)
	}

	return "", nil
}

package fze

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Runner(args []string) (string, error) {
	fmt.Printf("args: %v\n", args)
	argString := strings.Join(args, " ")
	switch cmd := args[0]; cmd {
	case "ls":
		return lsRunner(args[1:])
	case "find":
		return findRunner(args[1:])
	}
	out, err := exec.Command("bash", "-c", argString).Output()
	if err != nil {
		return "", fmt.Errorf("running command %v: %v", args, err)
	}
	return string(out), nil
}

func lsRunner(args []string) (string, error) {
	var cmd string
	if len(args) > 0 && args[0] == "-l" {
		cmd = "ls -l" + strings.Join(args, " ")
	} else {
		cmd = "ls " + strings.Join(args, " ")
	}
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func findRunner(args []string) (string, error) {
	// Get the output from find
	cmd := "find " + strings.Join(args, " ") // + " | fzf | xargs -n 1 emacsclient -n -s $TMUX_EMACS_DAEMON"
	fmt.Printf("Running cmd: %s\n", cmd)
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("running find: %v", err)
	}

	// Run the output from find through fzf
	fzf := exec.Command("fzf", "--tac")
	var out bytes.Buffer
	fzf.Stdin = bytes.NewReader(res)
	fzf.Stdout = &out
	fzf.Stderr = os.Stderr
	err = fzf.Run()
	if err != nil {
		return "", fmt.Errorf("starting fzf: %v", err)
	}

	path := strings.TrimSpace(out.String())
	fmt.Printf("running emacsclient on file %v", path)
	ec := exec.Command("emacsclient", "-n", "-s", "session3-window6", path)
	err = ec.Run()

	if err != nil {
		return "", fmt.Errorf("running emacsclient: %v", err)
	}

	return "", nil
}

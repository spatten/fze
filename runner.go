package fze

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
	case "rg":
		return rgRunner(args[1:])
	}
	out, err := exec.Command("bash", "-c", argString).Output()
	if err != nil {
		return "", fmt.Errorf("running command %v: %v", args, err)
	}
	return string(out), nil
}

func lsRunner(args []string) (string, error) {
	isLong := false
	if len(args) > 0 && args[0] == "-l" {
		isLong = true
	}

	cmd := "ls " + strings.Join(args, " ")
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	// Run the output from ls through fzf
	path, err := runFzf(res)
	if err != nil {
		return "", fmt.Errorf("runFzf: %v", err)
	}

	if isLong {
		re := regexp.MustCompile(" +")
		// ls -l output has a variable number of spaces between args. Clean this up by replacing multiple spaces with one space
		path = re.ReplaceAllString(path, " ")
		// take 9th column onwards as the path
		path = strings.SplitAfterN(path, " ", 9)[8]
		fmt.Printf("path arg: %s\n", path)
	}
	// Run emacsclient
	fmt.Printf("running emacsclient on file %v\n", path)
	ec := exec.Command("emacsclient", "-n", "-s", os.ExpandEnv("$TMUX_EMACS_DAEMON"), path)
	err = ec.Run()

	if err != nil {
		return "", fmt.Errorf("running emacsclient: %v", err)
	}

	return "", nil

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
	path, err := runFzf(res)
	if err != nil {
		return "", fmt.Errorf("runFzf: %v", err)
	}

	// Run emacsclient
	fmt.Printf("running emacsclient on file %v\n", path)
	ec := exec.Command("emacsclient", "-n", "-s", os.ExpandEnv("$TMUX_EMACS_DAEMON"), path)
	err = ec.Run()

	if err != nil {
		return "", fmt.Errorf("running emacsclient: %v", err)
	}

	return "", nil
}

func rgRunner(args []string) (string, error) {
	// Get the output from rg
	cmd := "rg -n " + strings.Join(args, " ") // + " | fzf | xargs -n 1 emacsclient -n -s $TMUX_EMACS_DAEMON"
	fmt.Printf("Running cmd: %s\n", cmd)
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("running rg: %v", err)
	}

	// Run the output from rg through fzf
	out, err := runFzf(res)
	if err != nil {
		return "", fmt.Errorf("runFzf: %v", err)
	}

	// Get the filename and linenumber from the output
	output := strings.Split(out, ":")
	if len(output) < 2 {
		return "", fmt.Errorf("expecting a path and a line-number in this rg output: %s", output)
	}
	path := output[0]
	lineNumber := output[1]

	// Run emacsclient

	fmt.Printf("running emacsclient on file %v\n", path)
	ec := exec.Command("emacsclient", "-n", "-s", os.ExpandEnv("$TMUX_EMACS_DAEMON"), "+"+lineNumber, path)
	err = ec.Run()

	if err != nil {
		return "", fmt.Errorf("running emacsclient with lineNumber = %s and path = %s: %v", lineNumber, path, err)
	}

	return "", nil
}

func runFzf(input []byte) (string, error) {
	fzf := exec.Command("fzf", "--tac")
	var out bytes.Buffer
	fzf.Stdin = bytes.NewReader(input)
	fzf.Stdout = &out
	fzf.Stderr = os.Stderr
	err := fzf.Run()
	if err != nil {
		return "", fmt.Errorf("starting fzf: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

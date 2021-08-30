package fze

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Runner(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no args provided")
	}

	var runnerOpts runnerOptions
	if args[0] == "-m" || args[0] == "--multi" {
		args = args[1:]
		runnerOpts = runnerOptions{multi: true}
	}

	switch cmd := args[0]; cmd {
	case "ls":
		return lsRunner(args[1:], runnerOpts)
	case "find":
		return findRunner(args[1:], runnerOpts)
	case "rg":
		return rgRunner(args[1:], runnerOpts)
	case "git":
		return gitRunner(args[1:], runnerOpts)
	case "gd": // alias for "git diff"
		return gitRunner(append([]string{"diff"}, args[1:]...), runnerOpts)
	case "st": // alias for "git status"
		return gitRunner(append([]string{"status"}, args[1:]...), runnerOpts)
	}
	return fmt.Errorf("command %s not recognized", args[0])
}

func openEditor(paths []pathArg, runnerOpts runnerOptions) error {
	var pathArgs []string

	// TODO: validate that no paths have a line-number or all paths have a line-number
	for _, path := range paths {
		if path.lineNumber != "" {
			fmt.Printf("line-number is %v\n", path.lineNumber)
			pathArgs = append(pathArgs, fmt.Sprintf("+%s", path.lineNumber))
		}
		pathArgs = append(pathArgs, path.path)
	}
	log.Printf("paths: %v (len = %d)\n", paths, len(paths))
	log.Printf("pathArgs: %v (len = %d)\n", pathArgs, len(pathArgs))
	var cmdArgs []string
	if !runnerOpts.multi {
		cmdArgs = []string{"-n"}
	}
	cmdArgs = append(cmdArgs, "-s", os.ExpandEnv("$TMUX_EMACS_DAEMON"))
	cmdArgs = append(cmdArgs, pathArgs...)
	ec := exec.Command("emacsclient", cmdArgs...)
	fmt.Printf("running command %v\n", ec)
	err := ec.Run()
	if err != nil {
		return fmt.Errorf("running emacsclient with args: %v, %g", cmdArgs, err)
	}

	tmux := exec.Command("tmux", "select-pane", "-U")
	err = tmux.Run()
	if err != nil {
		return fmt.Errorf("switching tmux pane: %v", err)
	}

	return nil
}

func runFzf(input []byte, opts runnerOptions) ([]string, error) {
	fzfArgs := []string{"--ansi", "--tac"}
	if opts.multi {
		fzfArgs = append(fzfArgs, "--multi")
	}
	fzf := exec.Command("fzf", fzfArgs...)
	var out bytes.Buffer
	fzf.Stdin = bytes.NewReader(input)
	fzf.Stdout = &out
	fzf.Stderr = os.Stderr
	err := fzf.Run()
	if err != nil {
		return nil, fmt.Errorf("starting fzf: %v", err)
	}

	var lines []string
	for _, line := range strings.Split(out.String(), "\n") {
		if line != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines, nil
}

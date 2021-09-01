package fze

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Runner(args []string, runnerOpts RunnerOptions) error {
	if len(args) < 1 {
		return fmt.Errorf("no args provided")
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

// fixArgs takes each arg, replaces ~ with the contents of $HOME, expands env variables and then
// wraps it in double-quotes.
// If then returns all of the fixed args separated by a space
func fixArgs(args []string) string {
	fixedArgs := make([]string, 0, len(args))
	home := os.ExpandEnv("$HOME")
	for _, arg := range args {
		fixed := "\"" + os.ExpandEnv(strings.ReplaceAll(arg, "~", home)) + "\""
		fixedArgs = append(fixedArgs, fixed)
	}
	return strings.Join(fixedArgs, " ")
}

func openEditor(paths []pathArg, runnerOpts RunnerOptions) error {
	var pathArgs []string

	// TODO: validate that no paths have a line-number or all paths have a line-number
	for _, path := range paths {
		if path.lineNumber != "" {
			pathArgs = append(pathArgs, fmt.Sprintf("+%s", path.lineNumber))
		}
		pathArgs = append(pathArgs, path.path)
	}
	var cmdArgs []string
	if !runnerOpts.Multi {
		cmdArgs = []string{"-n"}
	}
	cmdArgs = append(cmdArgs, "-s", runnerOpts.EmacsServer)
	cmdArgs = append(cmdArgs, pathArgs...)
	ec := exec.Command("emacsclient", cmdArgs...)
	err := ec.Run()
	if err != nil {
		return fmt.Errorf("running emacsclient with args: %v, %g", cmdArgs, err)
	}

	if runnerOpts.SwitchToTmuxPane {
		tmux := exec.Command("tmux", "select-pane", "-U")
		err = tmux.Run()
		if err != nil {
			return fmt.Errorf("switching tmux pane: %v", err)
		}
	}

	return nil
}

func runFzf(input []byte, opts RunnerOptions) ([]string, error) {
	fzfArgs := []string{"--ansi", "--tac"}
	if opts.Multi {
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

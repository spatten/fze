package fze

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Runner(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no args provided")
	}

	switch cmd := args[0]; cmd {
	case "ls":
		return lsRunner(args[1:])
	case "find":
		return findRunner(args[1:])
	case "rg":
		return rgRunner(args[1:])
	case "git":
		return gitRunner(args[1:])
	}
	return fmt.Errorf("command %s not recognized", args[0])
}

func openEditorWithLineNumber(path string, lineNumber string) error {
	lineNumberArg := fmt.Sprintf("+%s", lineNumber)
	return openEditor([]string{lineNumberArg, path})
}

func openEditorWithoutLineNumber(path string) error {
	return openEditor([]string{path})
}

func openEditor(pathArgs []string) error {
	args := append([]string{"-n", "-s", os.ExpandEnv("$TMUX_EMACS_DAEMON")}, pathArgs...)
	ec := exec.Command("emacsclient", args...)
	err := ec.Run()

	if err != nil {
		return fmt.Errorf("running emacsclient with args: %v, %g", pathArgs, err)
	}

	tmux := exec.Command("tmux", "select-pane", "-U")
	err = tmux.Run()
	if err != nil {
		return fmt.Errorf("switching tmux pane: %v", err)
	}

	return nil
}

func runFzf(input []byte) (string, error) {
	fzf := exec.Command("fzf", "--ansi", "--tac")
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

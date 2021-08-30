package fze

import (
	"fmt"
	"os/exec"
	"strings"
)

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
	err = openEditorWithoutLineNumber(path)
	if err != nil {
		return "", fmt.Errorf("running emacsclient: %v", err)
	}

	return "", nil
}

package fze

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func lsRunner(args []string) error {
	isLong := false
	if len(args) > 0 && args[0] == "-l" {
		isLong = true
	}

	cmd := "ls " + strings.Join(args, " ")
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return err
	}

	// Run the output from ls through fzf
	path, err := runFzf(res)
	if err != nil {
		return fmt.Errorf("runFzf: %v", err)
	}

	if isLong {
		re := regexp.MustCompile(" +")
		// ls -l output has a variable number of spaces between args. Clean this up by replacing multiple spaces with one space
		path = re.ReplaceAllString(path, " ")
		// take 9th column onwards as the path
		path = strings.SplitN(path, " ", 9)[8]
	}

	// Run emacsclient
	return openEditorWithoutLineNumber(path)
}

package fze

import (
	"fmt"
	"os/exec"
	"strings"
)

func gitRunner(args []string, opts RunnerOptions) error {
	if len(args) < 1 {
		return fmt.Errorf("no args to git")
	}

	switch arg := args[0]; arg {
	case "diff":
		return gitDiffRunner(args[1:], opts)
	case "show":
		return gitShowRunner(args[1:], opts)
	case "status":
		return gitStatusRunner(args[1:], opts)
	}

	return fmt.Errorf("only git diff and git show are supported")
}

func gitDiffRunner(args []string, opts RunnerOptions) error {
	// Get the output from git
	args, isStatus := mangleGitArgs(args)
	var cmd string
	if isStatus {
		cmd = "git diff --color=always " + fixArgs(args)
	} else {
		cmd = "git diff --src-prefix=a/ --dst-prefix=b/ --color=always " + fixArgs(args) + " | showlinenum.awk show_path=1"
	}

	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("running git command \"%s\": %v", cmd, err)
	}

	return gitDiffOrShowOutput(res, isStatus, opts)
}

func gitShowRunner(args []string, opts RunnerOptions) error {
	// Get the output from git
	args, isStatus := mangleGitArgs(args)
	var cmd string
	if isStatus {
		cmd = "git show --color=always --oneline" + fixArgs(args) + " | tail -n +2"
	} else {
		cmd = "git show --src-prefix=a/ --dst-prefix=b/ --color=always --oneline " + fixArgs(args) + " | tail -n +2 | showlinenum.awk show_path=1 "
	}

	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("running git command %v: %v", cmd, err)
	}

	return gitDiffOrShowOutput(res, isStatus, opts)
}

// replace "git status" with "git diff --stat"
func gitStatusRunner(args []string, opts RunnerOptions) error {
	args = append([]string{"--stat"}, args...)
	return gitDiffRunner(args, opts)
}

func mangleGitArgs(args []string) ([]string, bool) {
	newArgs := make([]string, 0, len(args))
	var isStatus bool
	for _, arg := range args {
		if arg == "--stat" || arg == "--numstat" {
			isStatus = true
			newArgs = append(newArgs, "--stat")
		} else {
			newArgs = append(newArgs, arg)
		}
	}
	return newArgs, isStatus
}

func gitDiffOrShowOutput(res []byte, isStatus bool, opts RunnerOptions) error {
	// Run the output from git through fzf
	outLines, err := runFzf(res, opts)
	if err != nil {
		return fmt.Errorf("runFzf: %v", err)
	}

	var paths []pathArg
	for _, line := range outLines {
		// git diff --stat has output like this:
		//
		// git diff --stat
		// gitRunner.go | 26 +++++++++++++++++++++++++++-----
		// lsRunner.go  |  2 +-
		// runner.go    |  4 ++++
		// 3 files changed, 66 insertions(+), 6 deletions(-)
		//
		// We want to grab everything from before the pipe
		if isStatus {
			path := strings.TrimSpace(strings.Split(line, "|")[0])
			paths = append(paths, pathArg{path: path})
		} else {
			// Get the filename and linenumber from the output
			output := strings.Split(line, ":")
			if len(output) < 2 {
				return fmt.Errorf("expecting a path and a line-number in this git output: %s", output)
			}
			paths = append(paths, pathArg{path: output[0], lineNumber: output[1]})
		}
	}

	// Run emacsclient
	err = openEditor(paths, opts)
	if err != nil {
		return fmt.Errorf("running emacsclient: %v", err)
	}

	return nil
}

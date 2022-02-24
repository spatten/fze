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
		return fzfRunner([]string{}, runnerOpts)
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
	case "ll": // alias for "ls -l"
		return lsRunner(append([]string{"-l"}, args[1:]...), runnerOpts)
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

	editorSetting := os.Getenv("FZE_EDITOR")
	if editorSetting == "" {
		editorSetting = "emacsclient"
	}

	var openEditorErr error
	switch editorSetting {
	case "vscode":
		openEditorErr = openVsCode(paths, runnerOpts)
	default:
		openEditorErr = openEmacsClient(paths, runnerOpts)
	}

	if openEditorErr != nil {
		return fmt.Errorf("opening editor %s: %v", editorSetting, openEditorErr)
	}

	return nil
}

func openEmacsClient(pathArgs []pathArg, runnerOpts RunnerOptions) error {
	var paths []string
	// TODO: validate that no paths have a line-number or all paths have a line-number
	for _, path := range pathArgs {
		if path.lineNumber != "" {
			paths = append(paths, fmt.Sprintf("+%s", path.lineNumber))
		}
		paths = append(paths, path.path)
	}
	var cmdArgs []string
	if !runnerOpts.Multi {
		cmdArgs = []string{"-n"}
	}
	cmdArgs = append(cmdArgs, "-s", runnerOpts.EmacsServer)
	cmdArgs = append(cmdArgs, paths...)
	ec := exec.Command("emacsclient", cmdArgs...)

	// Test mode
	if runnerOpts.TestFilter != "" {
		return nil
	}
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

func openVsCode(pathArgs []pathArg, runnerOpts RunnerOptions) error {
	var paths []string
	// TODO: validate that no paths have a line-number or all paths have a line-number
	for _, path := range pathArgs {
		if path.lineNumber != "" {
			paths = append(paths, fmt.Sprintf("%s:%s", path.path, path.lineNumber))
		}
		paths = append(paths, path.path)
	}
	var cmdArgs []string
	if runnerOpts.Multi {
		cmdArgs = []string{"--wait"}
	}
	cmdArgs = append(cmdArgs, "--goto")
	cmdArgs = append(cmdArgs, paths...)
	vsCode := exec.Command("code", cmdArgs...)

	// Test mode
	if runnerOpts.TestFilter != "" {
		return nil
	}
	err := vsCode.Run()
	if err != nil {
		return fmt.Errorf("running emacsclient with args: %v, %g", cmdArgs, err)
	}

	return nil
}

func runFzf(input []byte, opts RunnerOptions) ([]string, error) {
	fzfArgs := []string{"--ansi", "--tac"}
	if opts.Multi {
		fzfArgs = append(fzfArgs, "--multi")
	}

	// Used for testing. If you pass `--filter=foo` to fzf, then it
	// returns all lines that match "foo"
	if opts.TestFilter != "" {
		fzfArgs = append(fzfArgs, "--filter="+opts.TestFilter)
	}
	// fzf := exec.Command("fzf", fzfArgs...)
	fzf := exec.Command("fzf")
	var out bytes.Buffer
	fzf.Stdin = bytes.NewReader(input)
	fzf.Stdout = &out
	fzf.Stderr = os.Stderr
	err := fzf.Run()
	if err != nil {
		return nil, fmt.Errorf("starting fzf. input=%v\n Args = %v, err:: %v", string(input), fzfArgs, err)
	}

	var lines []string
	outLines := strings.Split(out.String(), "\n")

	// For testing -- just return the first match
	if opts.TestFilter != "" && !opts.Multi {
		outLines = []string{outLines[0]}
	}

	for _, line := range outLines {
		if line != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines, nil
}

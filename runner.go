package fze

import (
	"fmt"
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
	cmd := ""
	if len(args) == 0 {
		cmd = "ls"
	} else if args[0] == "-l" {
		cmd = "ls -l " + strings.Join(args[1:], " ")
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
	cmd := "find " + strings.Join(args, " ")
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(res), nil
}

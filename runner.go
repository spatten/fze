package fze

import (
	"fmt"
	"os/exec"
	"strings"
)

func Runner(args []string) (string, error) {
	fmt.Printf("args: %v\n", args)
	argString := strings.Join(args, " ")
	out, err := exec.Command("bash", "-c", argString).Output()
	if err != nil {
		return "", fmt.Errorf("running command %v: %v", args, err)
	}
	return string(out), nil
}

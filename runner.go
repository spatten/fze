package fze

import (
	"fmt"
	"os/exec"
)

func Runner(args []string) (string, error) {
	// out, err := exec.Command(args[0], args[1]).Output()
	fmt.Printf("args: %v\n", args)
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		return "", fmt.Errorf("running command %v: %v", args, err)
	}
	return string(out), nil
}

// Package smn_exec provides ...
package smn_exec

import (
	"os"
	"os/exec"
)

func EasyDirExec(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

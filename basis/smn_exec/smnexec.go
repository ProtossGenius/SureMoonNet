// Package smn_exec provides ...
package smn_exec

import (
	"bytes"
	"os"
	"os/exec"
)

//EasyDirExec run exec in dir.
func EasyDirExec(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

//DirExecGetOut get output.
func DirExecGetOut(dir, name string, args ...string) (oInfo, oErr string, err error) {
	cacheInfo := bytes.NewBuffer(nil)
	cacheErr := bytes.NewBuffer(nil)
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = cacheInfo
	cmd.Stderr = cacheErr
	err = cmd.Run()

	return cacheInfo.String(), cacheErr.String(), err
}

package utils

import "os/exec"

func Cmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}

func CmdWithOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	return string(output), err
}

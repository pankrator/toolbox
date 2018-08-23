package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Do(command string, args ...string) (string, error) {
	fmt.Println(command, args)
	cmd := exec.Command(command, args...)

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(stdErr.String())
		return "", err
	}

	return string(output), nil
}

func DoOutput(command string, args ...string) error {
	fmt.Println(command, args)
	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

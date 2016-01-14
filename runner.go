package genclient

import "os/exec"

type CommandRunner struct{}

func (c *CommandRunner) Run(cmd *exec.Cmd) error {
	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

package genclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

//go:generate counterfeiter --fake-name CommandRunner . CommandRunnerInterface
type CommandRunnerInterface interface {
	Run(cmd *exec.Cmd) error
}

type RPC struct {
	PathToBinary  string
	CommandRunner CommandRunnerInterface
}

func (r *RPC) ExecuteAndParse(methodName string, args map[string]interface{}, output interface{}) error {
	inputBytes, err := json.Marshal(map[string]interface{}{
		"Method": methodName,
		"Args":   args,
	})

	if err != nil {
		panic(err)
	}
	cmd := exec.Command(r.PathToBinary)
	cmd.Stdin = bytes.NewReader(inputBytes)
	stdoutBuffer, stderrBuffer := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer
	err = r.CommandRunner.Run(cmd)
	if err != nil {
		return err
	}

	err = json.Unmarshal(stdoutBuffer.Bytes(), &output)
	if err != nil {
		return fmt.Errorf("remote networker response cannot be parsed: %s\n%s\n%s",
			err, stdoutBuffer.Bytes(), stderrBuffer.Bytes())
	}
	return nil
}

package genclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/pivotal-golang/lager"
)

type ExternalNetworkerClient struct {
	path string
}

func New(path string) *ExternalNetworkerClient {
	return &ExternalNetworkerClient{path}
}

type methodCall struct {
	Method string
	Args   map[string]interface{} `json:",omitempty"`
}

func (e *ExternalNetworkerClient) Network(log lager.Logger, handle, spec string) (string, error) {
	cmd := exec.Command(e.path)
	call := methodCall{
		"Network",
		map[string]interface{}{
			"Handle": handle,
			"Spec":   spec,
		},
	}
	inputBytes, _ := json.Marshal(call)
	cmd.Stdin = bytes.NewReader(inputBytes)
	stdoutBuffer, stderrBuffer := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	ducatiErr := cmd.Wait()

	var output struct {
		Namespace string
		Error     string
	}
	err = json.Unmarshal(stdoutBuffer.Bytes(), &output)
	if err != nil {
		return "", fmt.Errorf("ducati response cannot be parsed: %s: %s", err, stdoutBuffer.Bytes())
	}

	if ducatiErr != nil {
		return "", fmt.Errorf("ducati failed: %s: %s", ducatiErr.Error(), output.Error)
	}
	return output.Namespace, nil
}

func (*ExternalNetworkerClient) Capacity() uint64 {
	// not implemented
	return 0
}

package genclient

import (
	"bytes"
	"encoding/json"
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
	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	var output struct {
		Namespace string
		Error     string
	}
	err = json.Unmarshal(stdoutBuffer.Bytes(), &output)
	if err != nil {
		panic(err)
	}
	return output.Namespace, nil
}

func (*ExternalNetworkerClient) Capacity() uint64 {
	// not implemented
	return 0
}

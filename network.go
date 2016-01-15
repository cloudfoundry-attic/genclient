package genclient

import (
	"errors"
	"fmt"
)

func (e *ExternalNetworkerClient) Network(handle, spec string) (string, error) {
	var output struct {
		Namespace string
		Error     string
	}
	err := e.RPC.ExecuteAndParse("Network", map[string]interface{}{
		"Handle": handle,
		"Spec":   spec,
	}, &output)

	if err != nil {
		return "", err
	}

	if output.Error != "" {
		return "", errors.New(output.Error)
	}

	if output.Namespace == "" {
		return "", fmt.Errorf("remote networker output missing Namespace")
	}
	return output.Namespace, nil
}

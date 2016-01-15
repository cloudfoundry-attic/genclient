package genclient

import "errors"

func (e *ExternalNetworkerClient) NetIn(handle string, hostPort int, containerPort int) error {
	var output struct {
		Error string
	}
	err := e.RPC.ExecuteAndParse("NetIn", map[string]interface{}{
		"Handle":        handle,
		"HostPort":      hostPort,
		"ContainerPort": containerPort,
	}, &output)

	if err != nil {
		return err
	}

	if output.Error != "" {
		return errors.New(output.Error)
	}

	return nil
}

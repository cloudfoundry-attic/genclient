package genclient

import "errors"

func (e *ExternalNetworkerClient) Destroy(handle string) error {
	var output struct {
		Error string
	}
	err := e.RPC.ExecuteAndParse("Destroy", map[string]interface{}{
		"Handle": handle,
	}, &output)

	if err != nil {
		return err
	}

	if output.Error != "" {
		return errors.New(output.Error)
	}

	return nil
}

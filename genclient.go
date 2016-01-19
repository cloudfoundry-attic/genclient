package genclient

//go:generate counterfeiter --fake-name RPC . RPCInterface
type RPCInterface interface {
	ExecuteAndParse(methodName string, args map[string]interface{}, output interface{}) error
}

type ExternalNetworkerClient struct {
	RPC RPCInterface
}

func New(ducatiBinaryPath, cniPluginDir string) *ExternalNetworkerClient {
	return &ExternalNetworkerClient{
		RPC: &RPC{
			PathToBinary:       ducatiBinaryPath,
			CommandRunner:      &CommandRunner{},
			CNIPluginDirectory: cniPluginDir,
		},
	}
}

func (*ExternalNetworkerClient) Capacity() uint64 {
	// not implemented
	return 0
}

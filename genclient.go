package genclient

//go:generate counterfeiter --fake-name RPC . RPCInterface
type RPCInterface interface {
	ExecuteAndParse(methodName string, args map[string]interface{}, output interface{}) error
}

type ExternalNetworkerClient struct {
	RPC RPCInterface
}

func New(path string) *ExternalNetworkerClient {
	return &ExternalNetworkerClient{
		RPC: &RPC{
			PathToBinary:  path,
			CommandRunner: &CommandRunner{},
		},
	}
}

func (*ExternalNetworkerClient) Capacity() uint64 {
	// not implemented
	return 0
}

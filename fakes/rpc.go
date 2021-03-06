// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/genclient"
)

type RPC struct {
	ExecuteAndParseStub        func(methodName string, args map[string]interface{}, output interface{}) error
	executeAndParseMutex       sync.RWMutex
	executeAndParseArgsForCall []struct {
		methodName string
		args       map[string]interface{}
		output     interface{}
	}
	executeAndParseReturns struct {
		result1 error
	}
}

func (fake *RPC) ExecuteAndParse(methodName string, args map[string]interface{}, output interface{}) error {
	fake.executeAndParseMutex.Lock()
	fake.executeAndParseArgsForCall = append(fake.executeAndParseArgsForCall, struct {
		methodName string
		args       map[string]interface{}
		output     interface{}
	}{methodName, args, output})
	fake.executeAndParseMutex.Unlock()
	if fake.ExecuteAndParseStub != nil {
		return fake.ExecuteAndParseStub(methodName, args, output)
	} else {
		return fake.executeAndParseReturns.result1
	}
}

func (fake *RPC) ExecuteAndParseCallCount() int {
	fake.executeAndParseMutex.RLock()
	defer fake.executeAndParseMutex.RUnlock()
	return len(fake.executeAndParseArgsForCall)
}

func (fake *RPC) ExecuteAndParseArgsForCall(i int) (string, map[string]interface{}, interface{}) {
	fake.executeAndParseMutex.RLock()
	defer fake.executeAndParseMutex.RUnlock()
	return fake.executeAndParseArgsForCall[i].methodName, fake.executeAndParseArgsForCall[i].args, fake.executeAndParseArgsForCall[i].output
}

func (fake *RPC) ExecuteAndParseReturns(result1 error) {
	fake.ExecuteAndParseStub = nil
	fake.executeAndParseReturns = struct {
		result1 error
	}{result1}
}

var _ genclient.RPCInterface = new(RPC)

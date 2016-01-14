package genclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGENClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Guardian External Networker Client Suite")
}

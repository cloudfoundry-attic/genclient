package genclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var pathToHappyFake, pathToSadFake string

func TestGENClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Guardian External Networker Client Suite")
}

var _ = BeforeSuite(func() {
	var err error
	pathToHappyFake, err = gexec.Build("github.com/cloudfoundry-incubator/genclient/fakes/happy")
	Expect(err).NotTo(HaveOccurred())

	pathToSadFake, err = gexec.Build("github.com/cloudfoundry-incubator/genclient/fakes/sad")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

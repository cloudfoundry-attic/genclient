package genclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"

	"github.com/cloudfoundry-incubator/genclient"
)

var _ = Describe("Guardian External Networker Client", func() {
	var (
		logger            lager.Logger
		externalNetworker *genclient.ExternalNetworkerClient
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		externalNetworker = genclient.New(pathToFake)
	})

	It("should forward the Network call to the external binary", func() {
		ns, err := externalNetworker.Network(logger, "some-handle", "some-spec")
		Expect(err).NotTo(HaveOccurred())

		Expect(ns).To(Equal("some-namespace"))
	})

	Context("when the external binary cannot be started", func() {
		It("should return an error", func() {
			externalNetworker = genclient.New("\t")
			_, err := externalNetworker.Network(logger, "some-handle", "some-spec")
			Expect(err).To(MatchError(ContainSubstring("executable file not found")))
		})
	})

	XContext("when the external binary exits with non-zero status code", func() {
		It("should return an error", func() {

		})
	})
})

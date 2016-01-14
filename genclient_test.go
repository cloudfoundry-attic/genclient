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
		externalNetworker = genclient.New(fakes.Binaries["happy"])
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

	Context("when the external binary exits with status 0 but the output is not parsable as JSON", func() {
		It("should return an error", func() {
			externalNetworker = genclient.New("echo")
			namespace, err := externalNetworker.Network(logger, "some-handle", "some-spec")
			Expect(err).To(MatchError("remote networker response cannot be parsed: unexpected end of JSON input: \n"))
			Expect(namespace).To(BeEmpty())
		})
	})

	Context("when the external binary exits with non-zero status code", func() {
		Context("when the output is parsable as JSON", func() {
			It("should return an error", func() {
				externalNetworker = genclient.New(fakes.Binaries["sad"])
				namespace, err := externalNetworker.Network(logger, "some-handle", "some-spec")
				Expect(err).To(MatchError("remote networker failed: exit status 17: something broke"))
				Expect(namespace).To(BeEmpty())
			})
		})
		Context("when the output is not parsable as JSON", func() {
			It("should return an error including the stdout and stderr", func() {
				externalNetworker = genclient.New(fakes.Binaries["troublesome"])
				namespace, err := externalNetworker.Network(logger, "some-handle", "some-spec")
				Expect(err).To(MatchError("remote networker failed: exit status 27:\nvery bad, no JSON\nsome log message"))
				Expect(namespace).To(BeEmpty())
			})
		})
	})

})

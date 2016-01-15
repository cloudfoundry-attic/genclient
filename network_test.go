package genclient_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/genclient"
	"github.com/cloudfoundry-incubator/genclient/fakes"
)

var _ = Describe("Network method", func() {
	var (
		rpc               *fakes.RPC
		externalNetworker *genclient.ExternalNetworkerClient
	)

	BeforeEach(func() {
		rpc = &fakes.RPC{}
		externalNetworker = &genclient.ExternalNetworkerClient{RPC: rpc}
	})

	It("should forward the Network call to the RPC", func() {
		rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
			Expect(methodName).To(Equal("Network"))
			json.Unmarshal([]byte(`{"Namespace": "some-namespace"}`), &output)

			Expect(args["Handle"]).To(Equal("some-handle"))
			Expect(args["Spec"]).To(Equal("some-spec"))
			return nil
		}

		ns, err := externalNetworker.Network("some-handle", "some-spec")
		Expect(err).NotTo(HaveOccurred())

		Expect(ns).To(Equal("some-namespace"))
	})

	Context("when RPC errors", func() {
		It("should return the error", func() {
			rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
				return errors.New("some error")
			}

			_, err := externalNetworker.Network("some-handle", "some-spec")
			Expect(err).To(MatchError("some error"))
		})
	})

	Context("when the result JSON does not include a Namespace", func() {
		It("should return an error", func() {
			rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
				json.Unmarshal([]byte(`{"Namespace": ""}`), &output)
				return nil
			}

			_, err := externalNetworker.Network("some-handle", "some-spec")
			Expect(err).To(MatchError(`remote networker output missing Namespace`))
		})
	})

	Context("when the resulting JSON includes an error string", func() {
		It("should return that string as an error, and not complain about missing namespace", func() {
			rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
				json.Unmarshal([]byte(`{"Error": "some error"}`), &output)
				return nil
			}

			ns, err := externalNetworker.Network("some-handle", "some-spec")
			Expect(err).To(MatchError("some error"))
			Expect(ns).To(BeEmpty())
		})

		It("should return that string as an error, and not return a namespace", func() {
			rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
				json.Unmarshal([]byte(`{"Namespace": "nonsense", "Error": "some error"}`), &output)
				return nil
			}

			ns, err := externalNetworker.Network("some-handle", "some-spec")
			Expect(err).To(MatchError("some error"))
			Expect(ns).To(BeEmpty())
		})
	})
})

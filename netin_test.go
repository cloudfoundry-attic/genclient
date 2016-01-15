package genclient_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/genclient"
	"github.com/cloudfoundry-incubator/genclient/fakes"
)

var _ = Describe("NetIn method", func() {
	var (
		rpc               *fakes.RPC
		externalNetworker *genclient.ExternalNetworkerClient
	)

	BeforeEach(func() {
		rpc = &fakes.RPC{}
		externalNetworker = &genclient.ExternalNetworkerClient{RPC: rpc}
	})

	It("should forward the NetIn call to the RPC", func() {
		rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
			Expect(methodName).To(Equal("NetIn"))

			Expect(args["Handle"]).To(Equal("some-handle"))
			Expect(args["HostPort"]).To(Equal(27777))
			Expect(args["ContainerPort"]).To(Equal(80))
			return nil
		}

		err := externalNetworker.NetIn("some-handle", 27777, 80)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when RPC errors", func() {
		It("should return the error", func() {
			rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
				return errors.New("some error")
			}

			err := externalNetworker.NetIn("some-handle", 27777, 80)
			Expect(err).To(MatchError("some error"))
		})
	})

	Context("when the resulting JSON includes an error string", func() {
		It("should return that string as an error", func() {
			rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
				json.Unmarshal([]byte(`{"Error": "some error"}`), &output)
				return nil
			}

			err := externalNetworker.NetIn("some-handle", 27777, 80)
			Expect(err).To(MatchError("some error"))
		})
	})
})

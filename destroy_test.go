package genclient_test

import (
	"encoding/json"
	"errors"

	"github.com/cloudfoundry-incubator/genclient"
	"github.com/cloudfoundry-incubator/genclient/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Destroy", func() {
	var (
		logger            lager.Logger
		rpc               *fakes.RPC
		externalNetworker *genclient.ExternalNetworkerClient
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		rpc = &fakes.RPC{}
		externalNetworker = &genclient.ExternalNetworkerClient{RPC: rpc}
	})

	It("should forward the Destroy call to the RPC", func() {
		rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
			Expect(methodName).To(Equal("Destroy"))

			Expect(args["Handle"]).To(Equal("some-handle"))
			return nil
		}

		err := externalNetworker.Destroy("some-handle")
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when RPC errors", func() {
		It("should return the error", func() {
			rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
				return errors.New("some error")
			}

			err := externalNetworker.Destroy("some-handle")
			Expect(err).To(MatchError("some error"))
		})
	})

	Context("when the resulting JSON includes an error string", func() {
		It("should return that string as an error", func() {
			rpc.ExecuteAndParseStub = func(methodName string, args map[string]interface{}, output interface{}) error {
				json.Unmarshal([]byte(`{"Error": "some error"}`), &output)
				return nil
			}

			err := externalNetworker.Destroy("some-handle")
			Expect(err).To(MatchError("some error"))
		})
	})
})

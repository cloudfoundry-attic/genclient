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

var _ = BeforeSuite(func() {
	fakes = newFakesRepo()
	Expect(fakes.Create("happy", `{ "Namespace": "some-namespace" }`, "", 0)).To(Succeed())
	Expect(fakes.Create("sad", `{ "Error": "something broke" }`, "", 17)).To(Succeed())
	Expect(fakes.Create("troublesome", `very bad, no JSON`, "some log message", 27)).To(Succeed())
})

var _ = AfterSuite(func() {
	fakes.Cleanup()
})

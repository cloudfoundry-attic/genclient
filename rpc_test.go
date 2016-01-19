package genclient_test

import (
	"errors"
	"io/ioutil"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/genclient"
	"github.com/cloudfoundry-incubator/genclient/fakes"
)

var _ = Describe("RPC", func() {
	var (
		rpc           *genclient.RPC
		commandRunner *fakes.CommandRunner
	)

	BeforeEach(func() {
		commandRunner = &fakes.CommandRunner{}

		rpc = &genclient.RPC{
			PathToBinary:       "/some/path/to/a/binary",
			CNIPluginDirectory: "/some/path/to/cni/bin",
			CommandRunner:      commandRunner,
		}
	})

	It("should marshal the args into JSON for stdin", func() {
		commandRunner.RunStub = func(cmd *exec.Cmd) error {
			stdinBytes, _ := ioutil.ReadAll(cmd.Stdin)
			Expect(stdinBytes).To(MatchJSON(` {
				"Method": "some-method-name",
				"Args": {
					"Arg1": "some-arg-1",
					"Arg2": 42
				}
			} `))
			cmd.Stdout.Write([]byte("{}"))
			return nil
		}

		var output struct{}
		err := rpc.ExecuteAndParse("some-method-name", map[string]interface{}{
			"Arg1": "some-arg-1",
			"Arg2": 42,
		}, &output)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should unmarshal stdout as JSON into an output struct", func() {
		var output struct {
			ReturnValue1 string
			ReturnValue2 int
			Error        string
		}
		commandRunner.RunStub = func(cmd *exec.Cmd) error {
			stdoutBytes := []byte(`{
				"ReturnValue1": "some-return-value-1",
				"ReturnValue2": -12345,
				"Error": "some-error-value"
			}`)
			cmd.Stdout.Write([]byte(stdoutBytes))
			return nil
		}

		err := rpc.ExecuteAndParse("some-method-name", map[string]interface{}{}, &output)
		Expect(err).NotTo(HaveOccurred())
		Expect(output.ReturnValue1).To(Equal("some-return-value-1"))
		Expect(output.ReturnValue2).To(Equal(-12345))
		Expect(output.Error).To(Equal("some-error-value"))
	})

	It("should set the CNI_PLUGIN_DIR env var on the ducati process", func() {
		commandRunner.RunStub = func(cmd *exec.Cmd) error {
			Expect(cmd.Env).To(ContainElement("CNI_PLUGIN_DIR=/some/path/to/cni/bin"))
			cmd.Stdout.Write([]byte("{}"))
			return nil
		}
		var output struct{}
		err := rpc.ExecuteAndParse("some-method-name", map[string]interface{}{}, &output)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when the command runner errors", func() {
		It("should return the error", func() {
			commandRunner.RunStub = func(cmd *exec.Cmd) error {
				return errors.New("some error")
			}

			var output struct{}
			err := rpc.ExecuteAndParse("some-method-name", map[string]interface{}{}, &output)
			Expect(err).To(MatchError("some error"))
		})
	})

	Context("when the output is not parsable as JSON", func() {
		It("should return an error including the stdout and stderr", func() {
			commandRunner.RunStub = func(cmd *exec.Cmd) error {
				cmd.Stdout.Write([]byte(`very bad, no JSON`))
				cmd.Stderr.Write([]byte(`some log message`))
				return nil
			}

			var output struct{}
			err := rpc.ExecuteAndParse("some-method-name", map[string]interface{}{}, &output)
			Expect(err).To(MatchError("remote networker response cannot be parsed: invalid character 'v' looking for beginning of value\nvery bad, no JSON\nsome log message"))
		})
	})
})

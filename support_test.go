package genclient_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/onsi/gomega/gexec"
)

var fakes *fakesRepo

type fakesRepo struct {
	gopath   string
	Binaries map[string]string
}

func newFakesRepo() *fakesRepo {
	gopath, err := ioutil.TempDir("", "temp-gopath")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Join(gopath, "src"), 0700)
	if err != nil {
		panic(err)
	}
	return &fakesRepo{
		gopath:   gopath,
		Binaries: make(map[string]string),
	}
}

func (f *fakesRepo) Create(name, stdout, stderr string, exitCode int) error {
	var err error
	f.Binaries[name], err = f.build(name, stdout, stderr, exitCode)
	return err
}

func (f *fakesRepo) Cleanup() {
	for _, v := range f.Binaries {
		os.Remove(v)
	}
	os.RemoveAll(f.gopath)
	gexec.CleanupBuildArtifacts()
}

func (f *fakesRepo) build(name, stdout, stderr string, exitCode int) (string, error) {
	srcCode := f.getSourceCode(stdout, stderr, exitCode)
	path := filepath.Join(f.gopath, "src", name, "main.go")
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path, []byte(srcCode), 0600)
	if err != nil {
		return "", err
	}

	return gexec.BuildIn(f.gopath, name)
}

func (*fakesRepo) getSourceCode(stdout, stderr string, exitCode int) string {
	template := `
package main

import (
	"os"
)

func main() {
	os.Stdout.Write([]byte(%q))
	os.Stderr.Write([]byte(%q))
	os.Exit(%d)
}`

	return fmt.Sprintf(template, stdout, stderr, exitCode)
}

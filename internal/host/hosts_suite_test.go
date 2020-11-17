package host_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHost(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "host state machine tests")
}

package conjungo

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMergeSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Merge Suite")
}

package conjungo

import (
	"testing"

	"github.com/Sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMergeSuite(t *testing.T) {
	// reduce the noise when testing
	logrus.SetLevel(logrus.FatalLevel)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Merge Suite")
}

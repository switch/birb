package failing_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFailing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Failing Suite")
}

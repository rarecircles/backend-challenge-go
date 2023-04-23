package storagerepo

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStoragerepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storagerepo Suite")
}

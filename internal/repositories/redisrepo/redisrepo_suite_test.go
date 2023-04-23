package redisrepo

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRedisrepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redisrepo Suite")
}

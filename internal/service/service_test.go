package service

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Service", func() {
	var s S

	BeforeEach(func() {
		s = New(zap.NewNop(), storageRepoMock{})
	})

	Describe("Searching", func() {
		Context("When searching for a word", func() {
			It("calls the repo with the same word", func() {
				query := "NETFLIX"
				result, err := s.Search(query)
				Expect(err).Should(MatchError(errors.New(query)))
				Expect(result).To(BeNil())
			})
		})

		Context("When searching for empty query", func() {
			It("calls the repo with empty query", func() {
				query := ""
				result, err := s.Search(query)
				Expect(err).Should(MatchError(errors.New(query)))
				Expect(result).To(BeNil())
			})
		})
	})
})

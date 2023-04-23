package storagerepo

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/internal/models"
	"go.uber.org/zap"
	"math/big"
)

var _ = Describe("Storagerepo", func() {
	var rr SR
	sampleAddress, _ := eth.NewAddress("0xbb5958767c8286bd1f86030e731549498e5323f7")
	sampleToken := models.Token{
		Name:         "Netflix",
		Symbol:       "NFLX",
		Address:      sampleAddress,
		Decimals:     18,
		TotalSupply:  big.NewInt(1000000000000000000),
		BaseTokenURI: "https://gateway.pin",
	}

	BeforeEach(func() {
		rr = New(zap.NewNop(), redisClientMock{})
	})

	Describe("Store", func() {
		Context("When storing a token in redis", func() {
			It("should store the correct properties", func() {
				err := rr.Store(sampleToken)
				Expect(err).Should(MatchError(errors.New(sampleToken.Name + sampleToken.Symbol)))
			})
		})
	})

	Describe("Search", func() {
		Context("When searching for empty", func() {
			It("should search for wildcard", func() {
				query := ""
				result, err := rr.Search(query)
				Expect(err).Should(MatchError(errors.New("*")))
				Expect(result).To(BeNil())
			})
		})

		Context("When searching for a query", func() {
			It("should search for the wildcard of the query", func() {
				query := "net"
				result, err := rr.Search(query)
				Expect(err).Should(MatchError(errors.New("*" + query + "*")))
				Expect(result).To(BeNil())
			})
		})
	})
})

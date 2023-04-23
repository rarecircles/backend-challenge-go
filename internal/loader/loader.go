package loader

import (
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/ethrepo"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/redisrepo"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/tokensrepo"
	"go.uber.org/zap"
	"time"
)

type Loader struct {
	ethRepo           ethrepo.IEthRepo
	redisRepo         IRedisRepo
	tokensRepo        tokensrepo.ITokensRepo
	l                 *zap.Logger
	refetchDelayHours int
	numWorkers        int
}

func New(l *zap.Logger, ethRepo ethrepo.IEthRepo, repo redisrepo.IRedisRepo, tokensRepo tokensrepo.ITokensRepo, numWorkers, refetchDelayHours int) Loader {
	loader := Loader{l: l, ethRepo: ethRepo, redisRepo: repo, tokensRepo: tokensRepo, numWorkers: numWorkers, refetchDelayHours: refetchDelayHours}
	return loader
}

func (l Loader) loadToken(semaphore <-chan bool, address eth.Address) {
	defer func() {
		<-semaphore
	}()
	for {
		token, err := l.ethRepo.GetToken(address)
		if err != nil {
			if err.Error()[len(err.Error())-3:] == "429" {
				l.l.Info("rate limit exceeded, waiting 1 second", zap.String("address", address.String()))
				time.Sleep(1 * time.Second)
				continue
			}
			l.l.Info("error getting token", zap.String("address", address.String()))
			return
		}
		if err = l.redisRepo.Store(token); err != nil {
			l.l.Error("error storing token", zap.Error(err))
			return
		}
		return
	}
}

func (l Loader) getDiffAddresses() ([]eth.Address, error) {
	allAddresses, err := l.tokensRepo.ListTokenAddresses()
	if err != nil {
		return nil, err
	}

	availableAddresses, _ := l.redisRepo.GetAllAddresses()

	var diffAddresses []eth.Address
	for _, address := range allAddresses {
		if !availableAddresses[address.String()] {
			diffAddresses = append(diffAddresses, address)
		}
	}
	return diffAddresses, nil
}

func (l Loader) loadAllTokens() error {
	tokenAddresses, err := l.getDiffAddresses()
	if err != nil {
		return err
	}

	//semaphore := make(chan bool, runtime.NumCPU())
	semaphore := make(chan bool, l.numWorkers)
	for _, address := range tokenAddresses {
		semaphore <- true
		go l.loadToken(semaphore, address)
	}
	return nil
}

func (l Loader) RunLoader() {
	go func() {
		for {
			err := l.loadAllTokens()
			if err != nil {
				l.l.Error("error loading tokens", zap.Error(err))
			}
			if l.refetchDelayHours == 0 {
				break
			}
			time.Sleep(time.Duration(l.refetchDelayHours) * time.Hour)
		}
	}()
}

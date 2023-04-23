package loader

import (
	"github.com/rarecircles/backend-challenge-go/eth"
	"github.com/rarecircles/backend-challenge-go/internal/models"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/ethrepo"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/storagerepo"
	"github.com/rarecircles/backend-challenge-go/internal/repositories/tokensrepo"
	"go.uber.org/zap"
	"time"
)

type L struct {
	ethRepo           ethrepo.I
	storageRepo       iStorageRepo
	tokensRepo        tokensrepo.I
	logger            *zap.Logger
	refetchDelayHours int
	numWorkers        int
}

func New(l *zap.Logger, ethRepo ethrepo.I, storageRepo storagerepo.I, tokensRepo tokensrepo.I, numWorkers, refetchDelayHours int) L {
	loader := L{logger: l, ethRepo: ethRepo, storageRepo: storageRepo, tokensRepo: tokensRepo, numWorkers: numWorkers, refetchDelayHours: refetchDelayHours}
	return loader
}

func (l L) loadToken(address eth.Address) (models.Token, error) {
	for {
		token, err := l.ethRepo.GetToken(address)
		if err != nil {
			if err.Error()[len(err.Error())-3:] == "429" {
				l.logger.Info("rate limit exceeded, waiting 1 second", zap.String("address", address.String()))
				time.Sleep(1 * time.Second)
				continue
			}
			return models.Token{}, err
		}
		return token, nil
	}
}

func (l L) storeToken(token models.Token) error {
	return l.storageRepo.Store(token)
}

func (l L) getDiffAddresses() ([]eth.Address, error) {
	allAddresses, err := l.tokensRepo.ListTokenAddresses()
	if err != nil {
		return nil, err
	}

	availableAddresses, _ := l.storageRepo.GetAllAddresses()

	var diffAddresses []eth.Address
	for _, address := range allAddresses {
		if !availableAddresses[address.String()] {
			diffAddresses = append(diffAddresses, address)
		}
	}
	return diffAddresses, nil
}

func (l L) loadAllTokens() error {
	tokenAddresses, err := l.getDiffAddresses()
	if err != nil {
		return err
	}

	// NOTE: limits the number of workers to avoid rate limit errors
	semaphore := make(chan bool, l.numWorkers)
	for _, address := range tokenAddresses {
		semaphore <- true
		go func(address eth.Address) {
			defer func() {
				<-semaphore
			}()
			token, err := l.loadToken(address)
			if err != nil {
				l.logger.Info("error getting token", zap.String("address", address.String()))
				return
			}
			err = l.storeToken(token)
			if err != nil {
				l.logger.Info("error storing token", zap.String("address", address.String()))
				return
			}
		}(address)
	}
	return nil
}

func (l L) RunLoader() {
	go func() {
		for {
			err := l.loadAllTokens()
			if err != nil {
				l.logger.Error("error loading tokens", zap.Error(err))
			}
			if l.refetchDelayHours == 0 {
				break
			}
			time.Sleep(time.Duration(l.refetchDelayHours) * time.Hour)
		}
	}()
}

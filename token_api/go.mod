module github.com/questionmarkquestionmark/go-go-backend/token_api

go 1.19

require (
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/go-test/deep v1.0.8
	github.com/gorilla/mux v1.8.0
	github.com/prometheus/client_golang v1.13.0
	github.com/rarecircles/backend-challenge-go v0.0.0-20211112223822-2839efbe5b80
	github.com/sethvargo/go-envconfig v0.8.2
	github.com/sirupsen/logrus v1.6.0
	github.com/slok/go-http-metrics v0.10.0
	github.com/tidwall/gjson v1.14.3
	go.uber.org/zap v1.23.0
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa
	golang.org/x/exp v0.0.0-20221002003631-540bb7301a08
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/rarecircles/backend-challenge-go/eth v0.0.0 => github.com/questionmarkquestionmark/go-go-backend/token_api/internal/eth v0.0.0

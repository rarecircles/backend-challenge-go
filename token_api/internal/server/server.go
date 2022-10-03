package server

// Implements the HTTP server for the token API.

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"

	"github.com/questionmarkquestionmark/go-go-backend/token_api/internal/api"
	"github.com/questionmarkquestionmark/go-go-backend/token_api/internal/dataloader"
)

// The Config struct provides all the configuration settings for the server.
type Config struct {
	// Hostname to use for the health checks and metrics.
	Host string

	// HTTP port to use.
	HTTPPort    int
	MetricsPort int

	// TLS file paths.
	TLSCert string `env:"TLS_CERT,default="`
	TLSKey  string `env:"TLS_KEY,default="`

	RPCURL        string `env:"RPC_URL,default=https://eth-mainnet.alchemyapi.io/v2/"`
	AlchemyAPIKey string `env:"API_URL,default=F6f7ymlA4Dr8fHjvVEsp1pxK-aCFq4LZ"`
	DataPath      string `env:"ADDRESS_DATA_PATH,default=./token_api/internal/dataloader/data/addresses.jsonl"`
}

// DefaultConfig builds up a default config that clients can use.
func DefaultConfig() Config {
	return Config{
		Host:        "0.0.0.0",
		HTTPPort:    8000,
		MetricsPort: 9000,
	}
}

func httpServer(c *Config) {
	log.Info("Beginning token data collection")
	data := dataloader.Process(c.RPCURL, c.AlchemyAPIKey, c.DataPath)
	log.Info("Completed token data collection")
	tokenMux, err := api.TokenMux(data)
	if err != nil {
		log.Fatalf("Error creating api.tokenMux. %s", err)
	}
	metricsMiddleWare := middleware.New(middleware.Config{
		Recorder: prometheus.NewRecorder(prometheus.Config{}),
	})
	tokenMux.Use(std.HandlerProvider("", metricsMiddleWare))

	host := fmt.Sprintf("%s:%d", c.Host, c.HTTPPort)

	//starting metrics server
	go func() {
		tokenHost := fmt.Sprintf("%s:%d", c.Host, c.MetricsPort)
		log.Infof("Running prometheus metric endpoint on %s", tokenHost)
		if !(c.TLSCert == "" || c.TLSKey == "") {
			if err := http.ListenAndServeTLS(tokenHost, c.TLSCert, c.TLSKey, promhttp.Handler()); err != nil {
				log.Fatalf("%s", err)
			}

		} else {
			log.Warning("TLS disabled as TLSCert and/or TLSKey are missing.")
			if err := http.ListenAndServe(tokenHost, promhttp.Handler()); err != nil {
				log.Fatalf("%s", err)
			}
		}

	}()

	if !(c.TLSCert == "" || c.TLSKey == "") {
		log.Infof("Running HTTPS Token API on %s", host)
		if err := http.ListenAndServeTLS(host, c.TLSCert, c.TLSKey, tokenMux); err != nil {
			log.Fatalf("%s", err)
		}

	} else {
		log.Warning("TLS disabled as TLSCert and/or TLSKey are missing.")
		log.Infof("Running HTTP Token API on %s", host)
		if err := http.ListenAndServe(host, tokenMux); err != nil {
			log.Fatalf("%s", err)
		}
	}
}

// Run launches the server.
func Run(c Config) {

	go httpServer(&c)

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
	log.Infof("Signal received, shutting down...")
}

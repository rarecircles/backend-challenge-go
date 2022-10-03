package main

import (
	"context"

	"github.com/sethvargo/go-envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/questionmarkquestionmark/go-go-backend/token_api/internal/server"
)

func main() {
	ctx := context.Background()

	config := server.DefaultConfig()
	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatalf("failed to process essential env vars: %s", err)
	}
	server.Run(config)
}

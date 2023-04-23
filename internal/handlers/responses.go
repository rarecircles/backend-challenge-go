package handlers

import "github.com/rarecircles/backend-challenge-go/internal/models"

type searchResponse struct {
	Tokens []models.Token `json:"tokens"`
}

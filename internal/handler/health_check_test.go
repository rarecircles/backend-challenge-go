package handler_test

import (
	"github.com/golang/mock/gomock"
	"github.com/jose-camilo/backend-challenge-go/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck_Handler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	c := e.NewContext(req, rec)
	h := &handler.Handlers{}

	err := h.HealthCheck(c)

	assert.Nil(t, err)
	expected := "{\"status\":\"OK\"}\n"
	assert.Equal(t, expected, rec.Body.String())
	assert.Equal(t, http.StatusOK, rec.Code)
}

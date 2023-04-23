package handlers

import "github.com/gin-gonic/gin"

type Interface interface {
	SearchTokens(c *gin.Context)
	Up(c *gin.Context)
}

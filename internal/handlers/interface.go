package handlers

import "github.com/gin-gonic/gin"

type I interface {
	SearchTokens(c *gin.Context)
	Up(c *gin.Context)
}

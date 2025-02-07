package primary

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	Login(c *gin.Context)
	Register (c *gin.Context)
}
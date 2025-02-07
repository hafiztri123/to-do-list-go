package primary

import "github.com/gin-gonic/gin"

type UserHandler interface {
	GetProfile(c *gin.Context)
}

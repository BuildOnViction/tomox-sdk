package endpoints

import "github.com/gin-gonic/gin"

func GinError(msg string) gin.H {
	return gin.H{
		"error": msg,
	}
}

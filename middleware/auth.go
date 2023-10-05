package middleware

import (
	"github.com/gin-gonic/gin"
	"invoice-view-be/utils/token"
	"net/http"
)

func JwtAuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Println(c.Cookie(''))
		err := token.IsTokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}

}

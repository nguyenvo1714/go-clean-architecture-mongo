package middleware

import (
	"github.com/gin-gonic/gin"
	"go-learning/clean-architecture-mongo/domain"
	"go-learning/clean-architecture-mongo/utils"
	"net/http"
	"strings"
)

func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		tokens := strings.Split(bearerToken, " ")
		if len(tokens) == 2 {
			token := tokens[1]
			authorized, err := utils.IsAuthorized(token, secret)
			if authorized {
				UserID, err := utils.ExtractIDFromToken(token, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
					c.Abort()
					return
				}

				c.Set("x-user-id", UserID)
				c.Next()
				return
			}

			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
			c.Abort()
			return
		}

		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized!"})
		c.Abort()
	}
}

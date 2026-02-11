package middleware

import (
	"chopper/internal/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwt *security.Jwt
}

func NewAuthMiddleware(jwt *security.Jwt) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func (a *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no token",
			})
			return
		}
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization field",
			})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := a.jwt.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
		c.Set("user_id", claims.Id)
		c.Set("username", claims.Username)
		c.Next()
	}
}

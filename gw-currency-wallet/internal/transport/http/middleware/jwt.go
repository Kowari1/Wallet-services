package middleware

import (
	"gw-currency-wallet/internal/models"
	"gw-currency-wallet/internal/pkg/logger"
	"gw-currency-wallet/internal/pkg/messages"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if len(tokenStr) < 8 || tokenStr[:7] != "Bearer " {
			logger.L.Warnf("unauthorized request: missing or malformed token, ip=%s", c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Error:   messages.MsgUnauthorized,
			})
			return
		}

		tokenStr = tokenStr[7:]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			logger.L.Warnf("invalid token: %v, ip=%s", err, c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Error:   messages.MsgInvalidToken,
				Details: err.Error(),
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.L.Error("failed to parse JWT claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Error:   messages.MsgInvalidToken,
			})
			return
		}

		if v, ok := claims["user_id"].(string); ok {
			c.Set("user_id", v)
		}

		c.Next()
	}
}

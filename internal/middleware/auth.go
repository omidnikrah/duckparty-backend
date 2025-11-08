package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/omidnikrah/duckparty-backend/internal/config"
)

type AuthUser struct {
	Email  string `json:"email"`
	UserID uint   `json:"userId"`
}

const AuthUserKey = "user"

func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue := c.GetHeader("Authorization")
		if tokenValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenValue = strings.TrimPrefix(tokenValue, "Bearer ")
		if tokenValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			email, _ := claims["email"].(string)
			subString, _ := claims["sub"].(string)

			sub, _ := strconv.ParseUint(subString, 10, 64)

			c.Set(AuthUserKey, AuthUser{
				Email:  email,
				UserID: uint(sub),
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

			c.Abort()
			return
		}
	}
}

func GetAuthUser(c *gin.Context) (AuthUser, bool) {
	value, exists := c.Get(AuthUserKey)
	if !exists {
		return AuthUser{}, false
	}

	user, ok := value.(AuthUser)
	if !ok {
		return AuthUser{}, false
	}

	return user, true
}

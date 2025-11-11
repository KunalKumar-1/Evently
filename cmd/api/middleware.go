package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func (app *application) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth header is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(app.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		userIdVal, exists := claims["UserId"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UserId not found in token"})
			c.Abort()
			return
		}

		userIdFloat, ok := userIdVal.(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId type in token"})
			c.Abort()
			return
		}
		user, err := app.models.Users.Get(int(userIdFloat))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access!"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

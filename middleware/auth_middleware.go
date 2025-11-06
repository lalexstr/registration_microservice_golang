package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"auth-service/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ContextUser struct {
	ID    uint
	Role  string
	Email string
}

const CtxUserKey = "currentUser"

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		// expected "Bearer <token>"
		var tokenString string
		_, err := fmt.Sscanf(auth, "Bearer %s", &tokenString)
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// validate alg
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenInvalidClaims
			}
			return config.JWTSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}
		// get sub (id)
		var uid uint64
		switch v := claims["sub"].(type) {
		case float64:
			uid = uint64(v)
		case string:
			u, _ := strconv.ParseUint(v, 10, 64)
			uid = u
		default:
			uid = 0
		}
		role := ""
		if r, ok := claims["role"].(string); ok {
			role = r
		}
		email := ""
		if e, ok := claims["email"].(string); ok {
			email = e
		}

		if uid == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid sub in token"})
			return
		}

		// set to context
		c.Set(CtxUserKey, &ContextUser{ID: uint(uid), Role: role, Email: email})
		c.Next()
	}
}

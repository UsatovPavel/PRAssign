package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"

	"github.com/UsatovPavel/PRAssign/internal/logging"
)

type ctxKey string

var ContextUserID ctxKey = "user_id"
var ContextIsAdmin ctxKey = "is_admin"

type Claims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func getToken(c *gin.Context) string {
	t := c.GetHeader("token")
	if t != "" {
		return t
	}
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

func AuthRequired(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			secret := viper.GetString("AUTH_KEY")
			if secret == "" {
				return nil, errors.New("auth key not configured")
			}
			return []byte(secret), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("is_admin", claims.IsAdmin)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ContextUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextIsAdmin, claims.IsAdmin)

		reqLogger := l.With(slog.String("user_id", claims.UserID), slog.Bool("is_admin", claims.IsAdmin))
		ctx = logging.WithLogger(ctx, reqLogger)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

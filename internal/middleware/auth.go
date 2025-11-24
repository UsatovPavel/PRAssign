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

func AuthRequired(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("token")
		if tokenString == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tokenString = strings.TrimPrefix(auth, "Bearer ")
			}
		}
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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

		var uid string
		var isAdmin bool

		if mc, ok := parsed.Claims.(jwt.MapClaims); ok {
			if v, ok := mc["user_id"].(string); ok {
				uid = v
			}
			switch v := mc["is_admin"].(type) {
			case bool:
				isAdmin = v
			case float64:
				isAdmin = v != 0
			case string:
				isAdmin = v == "true"
			default:
				isAdmin = false
			}
		} else {
			uid = claims.UserID
			isAdmin = claims.IsAdmin
		}

		if uid == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", uid)
		c.Set("is_admin", isAdmin)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ContextUserID, uid)
		ctx = context.WithValue(ctx, ContextIsAdmin, isAdmin)

		reqLogger := l.With(slog.String("user_id", uid), slog.Bool("is_admin", isAdmin))
		ctx = logging.WithLogger(ctx, reqLogger)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type TokenRequest struct {
	Username string `json:"username"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) TokenByUsername(c *gin.Context) {
	var req TokenRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	secret := viper.GetString("AUTH_KEY")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth key not configured"})
		return
	}
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"user_id":  req.Username,
		"is_admin": req.Username == "admin",
		"iat":      now.Unix(),
		"exp":      now.Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}
	c.JSON(http.StatusOK, TokenResponse{Token: signed})
}

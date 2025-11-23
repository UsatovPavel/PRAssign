package api

import (
	"github.com/gin-gonic/gin"
	
	"github.com/UsatovPavel/PRAssign/internal/api/health"
	"github.com/UsatovPavel/PRAssign/internal/api/team"
	"github.com/UsatovPavel/PRAssign/internal/api/users"
	"github.com/UsatovPavel/PRAssign/internal/api/pullrequests"
)

type Handlers struct {
	Team *team.Handler
	Users *users.Handler
	PR *pullrequest.Handler
}

func InitRouter(h *Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", health.HealthCheck)

	teamGroup := r.Group("/team")
	{
		teamGroup.POST("/add", h.Team.Add)
		teamGroup.GET("/get", h.Team.Get)
	}

	usersGroup := r.Group("/users")
	{
		usersGroup.POST("/setIsActive", h.Users.SetIsActive)
		usersGroup.GET("/getReview", h.Users.GetReview)
	}

	prGroup := r.Group("/pullRequest")
	{
		prGroup.POST("/create", h.PR.Create)
		prGroup.POST("/merge", h.PR.Merge)
		prGroup.POST("/reassign", h.PR.Reassign)
	}

	return r
}
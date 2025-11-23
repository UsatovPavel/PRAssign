package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/UsatovPavel/PRAssign/internal/api/health"
	"github.com/UsatovPavel/PRAssign/internal/api/pullrequest"
	"github.com/UsatovPavel/PRAssign/internal/api/team"
	"github.com/UsatovPavel/PRAssign/internal/api/users"
	"github.com/UsatovPavel/PRAssign/internal/middleware"
)

type Handlers struct {
	Team   *team.Handler
	Users  *users.Handler
	PR     *pullrequest.Handler
	Health *health.Handler
}

func InitRouter(h *Handlers, l *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(middleware.SkipK6Logger(l), gin.Recovery())
	r.GET("/health", h.Health.HealthCheck)

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

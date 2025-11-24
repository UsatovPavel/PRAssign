package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/UsatovPavel/PRAssign/internal/api/auth"
	"github.com/UsatovPavel/PRAssign/internal/api/health"
	"github.com/UsatovPavel/PRAssign/internal/api/pullrequest"
	"github.com/UsatovPavel/PRAssign/internal/api/statistics"
	"github.com/UsatovPavel/PRAssign/internal/api/team"
	"github.com/UsatovPavel/PRAssign/internal/api/users"
	"github.com/UsatovPavel/PRAssign/internal/middleware"
)

type Handlers struct {
	Team       *team.Handler
	Users      *users.Handler
	PR         *pullrequest.Handler
	Health     *health.Handler
	Auth       *auth.Handler
	Statistics *statistics.Handler
}

func InitRouter(h *Handlers, l *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(middleware.SkipK6Logger(l), gin.Recovery())
	r.POST("/auth/token", h.Auth.TokenByUsername)
	r.GET("/health", h.Health.HealthCheck)

	teamGroup := r.Group("/team")
	teamGroup.Use(middleware.AuthRequired(l))
	{
		teamGroup.POST("/add", h.Team.Add)
		teamGroup.GET("/get", h.Team.Get)
	}

	usersGroup := r.Group("/users")
	usersGroup.Use(middleware.AuthRequired(l))
	{
		usersGroup.POST("/setIsActive", h.Users.SetIsActive)
		usersGroup.GET("/getReview", h.Users.GetReview)
	}

	prGroup := r.Group("/pullRequest")
	prGroup.Use(middleware.AuthRequired(l))
	{
		prGroup.POST("/create", h.PR.Create)
		prGroup.POST("/merge", h.PR.Merge)
		prGroup.POST("/reassign", h.PR.Reassign)
	}

	statsGroup := r.Group("/statistics")
	statsGroup.Use(middleware.AuthRequired(l))
	{
		statsGroup.GET("/assignments/users", h.Statistics.AssignmentsByUsers)
		statsGroup.GET("/assignments/pullrequests", h.Statistics.AssignmentsByPullrequests)
		statsGroup.GET("/assignments/user/:id", h.Statistics.AssignmentsForUser)
	}

	return r
}

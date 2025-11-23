package main

import (
	"log/slog"
	"os"

	"github.com/UsatovPavel/PRAssign/internal/api"
	"github.com/UsatovPavel/PRAssign/internal/api/health"
	"github.com/UsatovPavel/PRAssign/internal/api/pullrequest"
	"github.com/UsatovPavel/PRAssign/internal/api/team"
	"github.com/UsatovPavel/PRAssign/internal/api/users"
	"github.com/UsatovPavel/PRAssign/internal/middleware"
	"github.com/UsatovPavel/PRAssign/internal/repository"
	"github.com/UsatovPavel/PRAssign/internal/service"
	"github.com/UsatovPavel/PRAssign/internal/storage"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	db, err := storage.NewPostgres()
	if err != nil {
		l.Error("db connection error", "err", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	prRepo := repository.NewPullRequestRepository(db)

	userService := service.NewUserService(userRepo, l)
	teamService := service.NewTeamService(teamRepo, l)
	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo, l)

	handlers := &api.Handlers{
		Team:   team.NewHandler(teamService, l),
		Users:  users.NewHandler(userService, l),
		PR:     pullrequest.NewHandler(prService, l),
		Health: health.NewHandler(l),
	}

	router := api.InitRouter(handlers, l)
	router.Use(middleware.SkipK6Logger(l))

	if err := router.Run(":8080"); err != nil {
		l.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

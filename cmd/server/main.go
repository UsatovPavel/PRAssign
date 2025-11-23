package main

import (
	"log"

	"github.com/UsatovPavel/PRAssign/internal/api"
	"github.com/UsatovPavel/PRAssign/internal/api/team"
	"github.com/UsatovPavel/PRAssign/internal/api/users"
	"github.com/UsatovPavel/PRAssign/internal/api/pullrequest"
	"github.com/UsatovPavel/PRAssign/internal/repository"
	"github.com/UsatovPavel/PRAssign/internal/service"
	"github.com/UsatovPavel/PRAssign/internal/storage"
)

func main() {
	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	prRepo := repository.NewPullRequestRepository(db)

	userService := service.NewUserService(userRepo)
	teamService := service.NewTeamService(teamRepo)
	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo)

	handlers := &api.Handlers{
		Team:  team.NewHandler(teamService),
		Users: users.NewHandler(userService),
		PR:    pullrequest.NewHandler(prService),
	}

	router := api.InitRouter(handlers)
	router.Run(":8080")
}

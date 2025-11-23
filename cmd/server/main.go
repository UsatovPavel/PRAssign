package main

import (
	"log"

	"github.com/gin-gonic/gin"
    "github.com/UsatovPavel/PRAssign/internal/api"
    "github.com/UsatovPavel/PRAssign/internal/storage"
    "github.com/UsatovPavel/PRAssign/internal/repository"
    "github.com/UsatovPavel/PRAssign/internal/service"
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
	prService := service.NewPRService(prRepo, teamRepo, userRepo)

	r := gin.Default()

	api.RegisterUserRoutes(r, userService)
	api.RegisterTeamRoutes(r, teamService)
	api.RegisterPullRequestRoutes(r, prService)

	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	r.Run(":8080")
}
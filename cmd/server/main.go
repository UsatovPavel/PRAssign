package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/viper"

	"github.com/UsatovPavel/PRAssign/internal/config"
	"github.com/UsatovPavel/PRAssign/internal/middleware"
	"github.com/UsatovPavel/PRAssign/internal/repository"
	"github.com/UsatovPavel/PRAssign/internal/service"
	"github.com/UsatovPavel/PRAssign/internal/storage"
	"github.com/UsatovPavel/PRAssign/internal/webapi"
	"github.com/UsatovPavel/PRAssign/internal/webapi/factorial"
	"github.com/UsatovPavel/PRAssign/internal/webapi/health"
	"github.com/UsatovPavel/PRAssign/internal/webapi/pullrequest"
	"github.com/UsatovPavel/PRAssign/internal/webapi/statistics"
	"github.com/UsatovPavel/PRAssign/internal/webapi/team"
	"github.com/UsatovPavel/PRAssign/internal/webapi/users"
)

func main() {
	l := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	if err := viper.BindEnv("AUTH_KEY"); err != nil {
		l.Error("viper bind env failed", "err", err)
		os.Exit(config.ExitCodeConfigError)
	}

	db, err := storage.NewPostgres()
	if err != nil {
		l.Error("db connection error", "err", err)
		os.Exit(config.ExitCodeConfigError)
	}

	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	prRepo := repository.NewPullRequestRepository(db)
	factorialRepo := repository.NewFactorialRepo(db.Pool)

	userService := service.NewUserService(userRepo, prRepo, l)
	teamService := service.NewTeamService(teamRepo, l)
	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo, l)

	factCfg, err := config.LoadFactorialKafkaConfig()
	if err != nil {
		l.Error("factorial kafka config", "err", err)
		os.Exit(config.ExitCodeConfigError)
	}
	factSvc, err := service.NewFactorialService(service.FactorialConfig{
		BootstrapServers: factCfg.Bootstrap,
		TopicTasks:       factCfg.TopicTasks,
	}, l)
	if err != nil {
		l.Error("factorial kafka producer init failed", "err", err)
		os.Exit(config.ExitCodeConfigError)
	}
	defer factSvc.Close()

	startResultsConsumer(factorialRepo, l)

	handlers := &webapi.Handlers{
		Team:       team.NewHandler(teamService, l),
		Users:      users.NewHandler(userService, l),
		PR:         pullrequest.NewHandler(prService, l),
		Health:     health.NewHandler(l),
		Statistics: statistics.NewHandler(prService, l),
		Factorial:  factorial.NewHandler(factSvc, factorialRepo),
	}

	router := webapi.InitRouter(handlers, l)
	router.Use(middleware.SkipK6Logger(l))

	if err := router.Run(":8080"); err != nil {
		l.Error("server stopped", "err", err)
		os.Exit(config.ExitCodeConfigError)
	}
}

func startResultsConsumer(factorialRepo repository.FactorialRepository, l *slog.Logger) {
	// Start factorial results consumer if enabled (default: on).
	if os.Getenv("FACTORIAL_RESULTS_CONSUMER_ENABLED") == "0" {
		return
	}
	resCfg, err := config.LoadFactorialResultsKafkaConfig()
	if err != nil {
		l.Error("factorial results kafka config", "err", err)
		os.Exit(config.ExitCodeConfigError)
	}
	consumer := service.NewFactorialResultConsumer(
		factorialRepo,
		service.FactorialResultConsumerConfig{
			Bootstrap: resCfg.Bootstrap,
			Group:     resCfg.Group,
			Topic:     resCfg.Topic,
		},
		l,
	)
	go func() {
		if err := consumer.Run(context.Background()); err != nil {
			l.Error("factorial result consumer stopped", "err", err)
			os.Exit(config.ExitCodeConfigError)
		}
	}()
}

package cmd

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/urfave/cli"
	"golang-skeleton/common/redis"
	"golang-skeleton/config"
	"golang-skeleton/docs" // docs is generated by Swag CLI, you have to import it.
	"golang-skeleton/domain/user"
	"golang-skeleton/handler"
	"golang-skeleton/repository"
	"golang-skeleton/router"
	"golang-skeleton/worker"
	"log"
)

// @title Todo Application
// @description This is a todo list management application
// @version 1.0
// @host localhost:8080
// @BasePath /v1
var ServerCMD = cli.Command{
	Name:    "server",
	Aliases: []string{"server"},
	Usage:   "Backend server",
	Action: func(ctx *cli.Context) (err error) {
		// ----- Init config
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatal("cannot load config: ", err)
		}

		// ----- Init logging
		//logger.Init()
		//logger.Client().Info("logger successfully\n")

		// ----- Init DB
		repoImpl := initDB(cfg)

		// ----- Init service
		userDomain := user.NewUserDomain(repoImpl)

		// ----- Init swagger
		initSwagger(cfg)

		// ----- Init worker
		worker.NewWorker(cfg)

		// ----- Init server handler
		r := router.New()
		r.GET("/swagger/*", echoSwagger.WrapHandler)
		v1 := r.Group("/v1")
		h := handler.NewHandler(cfg, userDomain)
		h.RegisterRoutes(v1)
		r.Logger.Fatal(r.Start(":" + cfg.ApiPort))
		return
	},
}

func initSwagger(cfg *config.Config) {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + cfg.ApiPort
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func initDB(cfg *config.Config) repository.IRepository {
	repository.InitGORM(cfg)
	repoImpl := repository.NewRepositoryImpl()
	return repoImpl
}

func initRedis(cfg *config.Config) redis.IService {
	return redis.NewRedisService(cfg)
}
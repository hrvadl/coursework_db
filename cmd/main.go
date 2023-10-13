package main

import (
	"go.uber.org/zap"

	"github.com/hrvadl/coursework_db/pkg/config"
	"github.com/hrvadl/coursework_db/pkg/controllers"
	"github.com/hrvadl/coursework_db/pkg/db"
	"github.com/hrvadl/coursework_db/pkg/repo"
	"github.com/hrvadl/coursework_db/pkg/server"
	"github.com/hrvadl/coursework_db/pkg/services"
	"github.com/hrvadl/coursework_db/pkg/templates"
)

func main() {
	logger := zap.Must(zap.NewDevelopment()).Sugar()
	logger.Info("Started initializing server...")

	logger.Info("Initializing the helper services...")
	cfg := config.New(logger).ParseFromEnv()
	cryptor := services.NewCryptor()
	tr := templates.NewResolver(cfg)

	logger.Info("Initializing the DB...")
	db := db.Must(db.New(cfg.DSN))

	logger.Info("Initializing the repositories...")
	stockRepo := repo.NewStock(db)
	emitentRepo := repo.NewEmitent(db)
	sessionRepo := repo.NewSession(db)
	dealRepo := repo.NewDeal(db)

	logger.Info("Initializing the core services...")
	stockService := services.NewStock(stockRepo, cryptor)
	dealService := services.NewDeal(dealRepo)
	emitentService := services.NewEmitent(emitentRepo, cryptor)
	authService := services.NewAuth(
		stockRepo,
		emitentRepo,
		sessionRepo,
		cryptor,
	)

	authController := controllers.NewAuth(authService, tr)
	profileController := controllers.NewProfile(emitentService, stockService, dealService, tr)
	dealController := controllers.NewDeal(dealService, tr)

	logger.Infof("Server is starting on port %v", cfg.ServerPort)
	srv := server.NewHTTP(&server.HTTPServerArgs{
		Session: sessionRepo,
		Logger:  logger,
		Controllers: &server.Controllers{
			Auth:    authController,
			Profile: profileController,
			Deal:    dealController,
		},
	})
	logger.Fatal(srv.ListenAndServe(cfg.ServerPort))
}

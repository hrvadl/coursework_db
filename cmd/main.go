package main

import (
	"go.uber.org/zap"

	"github.com/hrvadl/coursework_db/pkg/config"
	"github.com/hrvadl/coursework_db/pkg/controllers"
	"github.com/hrvadl/coursework_db/pkg/db"
	"github.com/hrvadl/coursework_db/pkg/repo"
	"github.com/hrvadl/coursework_db/pkg/server"
	"github.com/hrvadl/coursework_db/pkg/services"
)

func main() {
	logger := zap.Must(zap.NewDevelopment()).Sugar()
	logger.Info("Started initializing server...")

	logger.Info("Initializing the helper services...")
	cfg := config.New(logger).ParseFromEnv()
	cryptor := services.NewCryptor()

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
	authService := services.NewAuth(
		stockRepo,
		emitentRepo,
		sessionRepo,
		cryptor,
	)

	logger.Infof("Server is starting on port %v", cfg.ServerPort)
	srv := server.NewHTTP(&server.HTTPServerArgs{
		Logger: logger,
		Controllers: &server.Controllers{
			Auth: controllers.NewUser(authService),
		},
	})
	logger.Fatal(srv.ListenAndServe(cfg.ServerPort))
}

package main

import (
	"go.uber.org/zap"

	"github.com/hrvadl/coursework_db/pkg/config"
	"github.com/hrvadl/coursework_db/pkg/controllers"
	"github.com/hrvadl/coursework_db/pkg/db"
	"github.com/hrvadl/coursework_db/pkg/middleware"
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
	userRepo := repo.NewStock(db)
	sessionRepo := repo.NewSession(db)
	dealRepo := repo.NewDeal(db)
	securityRepo := repo.NewSecurity(db)
	inventoryRepo := repo.NewInventory(db)
	transactionRepo := repo.NewTransaction(db, inventoryRepo)

	logger.Info("Initializing the middlewares...")
	authM := middleware.NewAuth(sessionRepo, tr)
	logM := middleware.NewHTTPLogger(logger)
	corsM := middleware.NewCors()

	logger.Info("Initializing the core services...")
	inventoryService := services.NewInventory(inventoryRepo, userRepo)
	userService := services.NewStock(userRepo, cryptor)
	dealService := services.NewDeal(dealRepo, inventoryRepo, userRepo, transactionRepo)
	transactionService := services.NewTransaction(transactionRepo)
	securityService := services.NewSecurity(securityRepo)
	authService := services.NewAuth(
		userRepo,
		sessionRepo,
		cryptor,
	)

	authController := controllers.NewAuth(authService, tr)
	transactionController := controllers.NewTransaction(
		userService,
		dealService,
		inventoryService,
		tr,
	)
	profileController := controllers.NewProfile(
		userService,
		dealService,
		transactionService,
		securityService,
		tr,
	)
	dealController := controllers.NewDeal(
		dealService,
		securityService,
		inventoryService,
		tr,
	)
	inventoryController := controllers.NewInventory(
		inventoryService,
		userService,
		tr,
	)

	logger.Infof("Server is starting on port %v", cfg.ServerPort)
	srv := server.NewHTTP(&server.HTTPServerArgs{
		AuthM:   authM,
		CorsM:   corsM,
		LoggerM: logM,
		Session: sessionRepo,
		Logger:  logger,
		Controllers: &server.Controllers{
			Auth:        authController,
			Profile:     profileController,
			Deal:        dealController,
			Inventory:   inventoryController,
			Transaction: transactionController,
		},
	})
	logger.Fatal(srv.ListenAndServe(cfg.ServerPort))
}

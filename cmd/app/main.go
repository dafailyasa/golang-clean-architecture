package main

import (
	"auth-service/config"
	"auth-service/delivery/http/router"
	"auth-service/infrastructure/persistence/mysql"
	pkgAppHttp "auth-service/pkg/app_http"
	pkgLogger "auth-service/pkg/logger"
	pkgValidator "auth-service/pkg/validator"
)

func main() {
	logger := pkgLogger.NewLogrusLogger()

	cfg, err := config.LoadConfigPath()
	if err != nil {
		logger.Fatal(err)
	}

	validator, ut, err := pkgValidator.NewValidator()
	if err != nil {
		logger.Fatal(err)
	}

	db, err := mysql.NewGormDB(&cfg.Database)
	if err != nil {
		logger.Fatal(err)
	}

	txmManager, err := mysql.NewTransactionManager(db)
	if err != nil {
		logger.Fatal(err)
	}

	httpClient := pkgAppHttp.NewClient(logger)

	muxRouter := router.NewRouter(cfg, txmManager, db, validator, httpClient, ut)
	if err := router.RunHTTPServer(muxRouter, cfg.Server.Port, logger); err != nil {
		logger.Fatalf("Failed to create http router: %v", err)
	}
}

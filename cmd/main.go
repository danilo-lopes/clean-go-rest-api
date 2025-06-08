// Clean Architecture - Main entrypoint
package main

import (
	"database/sql"
	"fmt"

	"clean-go-rest-api/internal/adapter/handler"
	"clean-go-rest-api/internal/adapter/repository"
	"clean-go-rest-api/internal/config"
	"clean-go-rest-api/internal/crosscutting/logger"
	"clean-go-rest-api/internal/infrastructure/db"
	"clean-go-rest-api/internal/infrastructure/server"
	"clean-go-rest-api/internal/usecase"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	logger := logger.NewLogger()

	logger.Info("Starting API application")

	logger.Info("Loading configuration")
	cfg := config.LoadConfig()

	logger.Info("Running migrations")
	dbParametersDefault := fmt.Sprintf("sslmode=disable TimeZone=%s", cfg.TimeZone)
	cfg.DB.Parameters = dbParametersDefault + " application_name=migrations"
	migrationsDBConn, err := sql.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to the database: %s", err.Error()))
		panic(err)
	}
	if err := migrationsDBConn.Ping(); err != nil {
		logger.Error(fmt.Sprintf("Error when trying to verify the connection to the database: %s", err.Error()))
		panic(err)
	}
	migration := db.NewMigration(migrationsDBConn, cfg.DB.MigrationsFolderPath)
	migration.RunMigrations()
	logger.Info("Finished migrations script")

	cfg.DB.Parameters = dbParametersDefault + " application_name=go_rest_api"
	dbConn, err := sql.Open("postgres", cfg.DBConnectionString())
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to connect to the database: %s", err.Error()))
		panic(err)
	}
	if err := dbConn.Ping(); err != nil {
		logger.Error(fmt.Sprintf("Error when trying to verify the connection to the database: %s", err.Error()))
		panic(err)
	}
	repo := repository.NewPostgresUserRepository(dbConn)
	userUseCase := usecase.NewUserUseCase(repo)
	router := mux.NewRouter()
	handler.NewUserHandler(userUseCase, logger).RegisterRoutes(router)
	server.StartServer(router, cfg.ServerPort, logger)
}

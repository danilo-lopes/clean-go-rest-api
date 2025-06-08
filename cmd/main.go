// Clean Architecture - Main entrypoint
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func initLogger() logger.ILogger {
	return logger.NewLogger()
}

func loadConfig(logger logger.ILogger) *config.Config {
	logger.Info("Loading configuration")
	cfg := config.LoadConfig()
	return cfg
}

func runMigrations(cfg *config.Config, logger logger.ILogger) {
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
}

func initDB(cfg *config.Config, logger logger.ILogger) *sql.DB {
	dbParametersDefault := fmt.Sprintf("sslmode=disable TimeZone=%s", cfg.TimeZone)
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
	return dbConn
}

func setupRouter(dbConn *sql.DB, logger logger.ILogger) *mux.Router {
	repo := repository.NewPostgresUserRepository(dbConn)
	userUseCase := usecase.NewUserUseCase(repo)
	router := mux.NewRouter()
	handler.NewUserHandler(userUseCase, logger).RegisterRoutes(router)
	handler.NewHealthCheckHandler(dbConn).RegisterRoutes(router)
	return router
}

func startServer(router *mux.Router, port int, logger logger.ILogger) *http.Server {
	httpServer := server.StartServer(router, port, logger)
	logger.Info(fmt.Sprintf("Server running on port %d", port))
	return httpServer
}

func main() {
	logger := initLogger()
	logger.Info("Starting API application")

	cfg := loadConfig(logger)

	runMigrations(cfg, logger)

	dbConn := initDB(cfg, logger)
	router := setupRouter(dbConn, logger)
	httpServer := startServer(router, cfg.ServerPort, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	} else {
		logger.Info("Server exited gracefully")
	}
}

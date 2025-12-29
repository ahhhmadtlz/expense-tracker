package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/ahhhmadtlz/expense-tracker/internal/config"
	"github.com/ahhhmadtlz/expense-tracker/internal/delivery/httpserver"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	userrepository "github.com/ahhhmadtlz/expense-tracker/internal/domain/user/repository"
	userservice "github.com/ahhhmadtlz/expense-tracker/internal/domain/user/service"
	uservalidator "github.com/ahhhmadtlz/expense-tracker/internal/domain/user/validator"
	categoryrepository "github.com/ahhhmadtlz/expense-tracker/internal/domain/category/repository"
	categoryservice "github.com/ahhhmadtlz/expense-tracker/internal/domain/category/service"
	categoryvalidator "github.com/ahhhmadtlz/expense-tracker/internal/domain/category/validator"
	transactionrepository "github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/repository"
	transactionservice "github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/service"
	transactionvalidator "github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/validator"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/repository/migrator"
	"github.com/ahhhmadtlz/expense-tracker/internal/repository/mysql"
)

func main() {
	cfg := config.C()

	appLogger := logger.New(
		cfg.Logger,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)
	// Create logs directory
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("❌ Failed to create logs directory: %v", err)
	}

	logger.SetDefault(appLogger)

	mysqlDB := mysql.New(cfg.Mysql)
	log.Println("✓ Database connected")

	if err := runMigrations(mysqlDB, cfg); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	appLogger.Info("Application starting", "port", cfg.HTTPServer.Port)

	authSvc, userSvc, userValidator, categorySvc, categoryValidator, transactionSvc, transactionValidator := setupServices(cfg, mysqlDB)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, categorySvc, categoryValidator, transactionSvc, transactionValidator)

	go func() {
		appLogger.Info("Server starting", "port", cfg.HTTPServer.Port)
		server.Serve()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	appLogger.Info("Received interrupt signal, shutting down...")

	// Close database connection
	if err := mysqlDB.Close(); err != nil {
		appLogger.Error("Failed to close database connection", "error", err.Error())
	} else {
		appLogger.Info("Database connection closed")
	}

	appLogger.Info("Application stopped")
}

func runMigrations(mysqlDB *mysql.MySQLDB, cfg config.Config) error {
	mgr := migrator.New(mysqlDB, migrator.Config{
		MigrationsDir:  "./internal/repository/mysql/migrations",
		MigrationTable: "schema_migrations",
	})

	return mgr.Up(0)
}


func setupServices(
	cfg config.Config,
	mysqlDB *mysql.MySQLDB,
) (
	auth.Service,
	userservice.Service,
	uservalidator.Validator,
	categoryservice.Service,
	categoryvalidator.Validator,
	transactionservice.Service,
	transactionvalidator.Validator,
) {
	// Auth service
	authSvc := auth.New(cfg.Auth)

	userRepo := userrepository.New(mysqlDB)
	userValidator := uservalidator.New(userRepo)
	userSvc := userservice.New(authSvc, userRepo)

	categoryRepo := categoryrepository.New(mysqlDB)
	categoryValidator := categoryvalidator.New(categoryRepo)
	categorySvc := categoryservice.New(categoryRepo)

	transactionRepo := transactionrepository.New(mysqlDB)
	transactionValidator := transactionvalidator.New(transactionRepo, categoryRepo)
	transactionSvc := transactionservice.New(transactionRepo)

	return authSvc, userSvc, userValidator, categorySvc, categoryValidator, transactionSvc, transactionValidator
}
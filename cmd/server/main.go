package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"vasek/internal/dto"
	"vasek/internal/handlers"
	"vasek/internal/repositories"
	"vasek/internal/server"
	"vasek/internal/services"
)

func main() {
	cfg := dto.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5434),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "vasek"),
		SSLMode:  getEnv("SSL_MODE", "disable"),
	}

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	err = runMigrations(db)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Could not run migrations: %v", err)
	}
	fmt.Println("Migrations ran successfully")

	repoService := repositories.NewRequestRepository(db)
	requestService := services.NewRequestService(db, *repoService)
	requestController := handlers.NewRequestHandler(requestService)

	newServer := server.NewServer(requestController)
	newServer.Start()
}

func connectDB(cfg dto.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations up: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	// Get database connection string from environment or use default
	// Priority: DATABASE_URL env var > default DSN
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5500/english_coach?sslmode=disable"
	}

	// Use pgx/stdlib for compatibility with sql.Open
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Read migration file
	migrationPath := "db/migrations/schema/0001_init_schema.sql"
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		// Try alternative path
		migrationPath = filepath.Join("backend", migrationPath)
	}

	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	// Execute migration
	sql := string(sqlBytes)
	if _, err := db.ExecContext(ctx, sql); err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	fmt.Println("Migration completed successfully!")
}

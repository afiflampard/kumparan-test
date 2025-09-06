package infra

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/afif-musyayyidin/hertz-boilerplate/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func InitPostgres(cfg config.Config) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open Postgres: %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping Postgres: %v", err)
	}

	log.Println("✅ Connected to Postgres")
	return db
}

func InitPostgresReplica(cfg config.Config) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBReplicaHost,
		cfg.DBReplicaPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open Postgres: %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping Postgres: %v", err)
	}

	log.Println("✅ Connected to Postgres Replica")
	return db
}

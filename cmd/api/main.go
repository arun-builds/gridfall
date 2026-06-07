package main

import (
	"context"
	"log"
	"os"

	"github.com/arun-builds/gridfall/internal/database"
	"github.com/arun-builds/gridfall/internal/database/db"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	ctx := context.Background()

	pool, err := database.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(pool)
	_ = queries
}

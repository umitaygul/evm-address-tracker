package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/umitaygul/evm-address-tracker/internal/api"
	"github.com/umitaygul/evm-address-tracker/internal/db"
)

func main() {
	// Local dev: .env varsa yükle, yoksa sistem env kullan
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	r := api.NewRouter(pool)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("api listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

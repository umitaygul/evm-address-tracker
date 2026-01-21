package main

import (
	"context"
	"log"

	"github.com/umitaygul/evm-address-tracker/internal/db"
)

func main() {
	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	log.Println("database connected")
}

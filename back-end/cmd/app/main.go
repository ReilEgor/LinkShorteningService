package main

import (
	"log"

	"github.com/ReilEgor/LinkShorteningService/internal/server"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("Critical error: %v", err)
	}
}

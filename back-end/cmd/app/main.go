package main

import (
	"log"

	"github.com/ReilEgor/CleanArchitectureGolang/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("Critical error: %v", err)
	}
}

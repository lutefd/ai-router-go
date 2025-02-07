package main

import (
	"log"

	"github.com/lutefd/ai-router-go/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

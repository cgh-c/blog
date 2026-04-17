package main

import (
	"log"
	"os"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Blog backend starting on port %s", port)
}

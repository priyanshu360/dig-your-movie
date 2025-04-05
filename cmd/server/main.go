package main

import (
	"log"

	"github.com/priyanshu360/dig-your-movie/internal/config"
	"github.com/priyanshu360/dig-your-movie/internal/dns"
	"github.com/priyanshu360/dig-your-movie/internal/omdb"
)

func main() {
	cfg := config.Load()

	omdbClient := omdb.NewClient(cfg.APIKey)
	server := dns.NewServer(cfg, omdbClient)

	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

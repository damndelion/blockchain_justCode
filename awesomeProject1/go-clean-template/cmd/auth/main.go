package main

import (
	"github.com/evrone/go-clean-template/config/auth"
	"github.com/evrone/go-clean-template/internal/auth/applicator"
	"log"
)

func main() {
	cfg, err := auth.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	applicator.Run(cfg)

}

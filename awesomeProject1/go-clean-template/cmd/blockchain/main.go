package main

import (
	blockchain2 "github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/blockchain"
	"log"
)

func main() {
	cfg, err := blockchain2.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	blockchain.Run(cfg)

}

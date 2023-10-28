package main

import (
	"github.com/evrone/go-clean-template/config/blockchain"
	"github.com/evrone/go-clean-template/internal/blockchain/applicator"
	"log"
)

func main() {
	cfg, err := blockchain.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	applicator.Run(cfg)

}

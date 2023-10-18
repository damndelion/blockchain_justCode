package main

import (
	"github.com/evrone/go-clean-template/config/user"
	"github.com/evrone/go-clean-template/internal/app"
	"log"
)

func main() {
	cfg, err := user.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)

}

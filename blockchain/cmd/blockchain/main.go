package main

import (
	"log"

	"github.com/damndelion/blockchain_justCode/config/blockchain"
	"github.com/damndelion/blockchain_justCode/internal/blockchain/applicator"
)

// @title Blockchain service
// @version 1.0
// @description Service that handles all blockchain request
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http
func main() {
	cfg, err := blockchain.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	applicator.Run(cfg)
}

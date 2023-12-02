package main

import (
	"log"

	"github.com/damndelion/blockchain_justCode/config/user"
	"github.com/damndelion/blockchain_justCode/internal/user/applicator"
)

// @title User service
// @version 1.0
// @description Service that does CRUD operations on user
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http.
func main() {
	cfg, err := user.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	applicator.Run(cfg)
}

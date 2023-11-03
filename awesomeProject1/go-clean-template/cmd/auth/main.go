package main

import (
	"github.com/evrone/go-clean-template/config/auth"
	"github.com/evrone/go-clean-template/internal/auth/applicator"
	"log"
)

// @title Authorization service
// @version 1.0
// @description Service that handles authorization and authentication and generates jwt token
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /
// @schemes http
func main() {
	cfg, err := auth.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	applicator.Run(cfg)

}

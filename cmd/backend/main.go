package main

import (
	"log"

	"github.com/antalkon/Go_prod_tmpl/internal/app"
)

// @title           Go Echo Template API
// @version         1.0
// @description     Go echo template API swagger documentation

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

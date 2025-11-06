package main

import (
	"log"

	"auth-service/config"
	"auth-service/db"
	"auth-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// init DB
	db.Init()

	r := gin.Default()

	// basic middleware - logging, recovery already included
	routes.Setup(r)

	addr := ":" + config.Port
	log.Printf("Starting auth service on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

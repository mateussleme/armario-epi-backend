package main

import (
	"log"

	"github.com/TopSisErp/epi-backend/internal/database"
	v1 "github.com/TopSisErp/epi-backend/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	err := database.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	router := gin.Default()
	router.Use(cors.Default())

	v1.Routes(router)

	router.Run("0.0.0.0:8080")
}

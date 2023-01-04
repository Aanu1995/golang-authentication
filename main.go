package main

import (
	"net/http"
	"os"

	"github.com/Aanu1995/golang-authentication/routes"
	"github.com/gin-gonic/gin"
)


func main(){
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	router := gin.Default()

	// Group all routes
	apiRouter := router.Group("/api")
	// Authentication routes
	routes.AuthRoute(apiRouter)
	// User routes
	routes.UserRoute(apiRouter)

	// default route
	apiRouter.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message":"Welcome to API version 1"})
	})

	router.Run(":" + port)
}
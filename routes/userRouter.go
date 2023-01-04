package routes

import (
	"github.com/Aanu1995/golang-authentication/controllers"
	"github.com/Aanu1995/golang-authentication/middleware"
	"github.com/gin-gonic/gin"
)


func UserRoute(router *gin.RouterGroup){
	// User routes
	userRouter := router.Group("/users", middleware.Authenticate)

	userRouter.GET("", controllers.GetUsers)
	userRouter.GET("/:userId", controllers.GetUser)
}
package middleware

import (
	"net/http"

	"github.com/Aanu1995/golang-authentication/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate(ctx *gin.Context){
	clientToken := ctx.GetHeader("token")
	if clientToken == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization header required"})
		ctx.Abort()
		return
	}

	claims, err := helpers.ValidateToken(clientToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Set("email", claims.Email)
	ctx.Set("firstName", claims.FirstName)
	ctx.Set("lastName", claims.LastName)
	ctx.Set("userId", claims.UserId)
	ctx.Set("userType", claims.UserType)
	ctx.Next()
}
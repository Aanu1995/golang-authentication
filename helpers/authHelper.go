package helpers

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func MatchUserTypeToUserID(ctx *gin.Context, userId string) (err error) {
	userType := ctx.GetString("userType")
	uid := ctx.GetString("userId")

	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access this resource")
	}
	return
}

func CheckUserType(ctx *gin.Context, role string) (err error) {
	userType := ctx.GetString("userType")

	if userType != role {
		err = errors.New("Unauthorized to access this resource")
	}
	return
}

// verify password
func VerifyPassword(userPassword string, providedPassword string) (passwordIsValid bool) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	passwordIsValid = err == nil
	return
}

// HashPassword
func Hashpassword(password string) string {
	result, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(result)
}
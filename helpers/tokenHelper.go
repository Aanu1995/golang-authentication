package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/Aanu1995/golang-authentication/models"
	"github.com/golang-jwt/jwt/v4"
)

type SignedDetails struct {
	Email 				string
	FirstName 		string
	LastName 			string
	UserId				string
	UserType			string
	jwt.RegisteredClaims
}

var secretkey = os.Getenv("SECRET_KEY")

func GenerateTokens(user models.User) (accessToken string, refreshToken string, err error){
	claims := SignedDetails{
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		UserId: user.UserId,
		UserType: user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	refreshClaims := SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(168 * time.Hour)),
		},
	}

	if accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretkey)); err != nil {
		return
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretkey))
	return
}


func ValidateToken(signedToken string) (claims *SignedDetails, err error){
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})

	if err != nil {
		return
	}


	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		err  = errors.New("Invalid authorization token")
		return
	}

	if claims.ExpiresAt.Before(time.Now()){
		err  = errors.New("Token is Expired")
		return
	}

	return
}
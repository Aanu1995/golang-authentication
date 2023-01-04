package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Aanu1995/golang-authentication/helpers"
	"github.com/Aanu1995/golang-authentication/models"
	"github.com/Aanu1995/golang-authentication/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func SignUp(ctx *gin.Context){
	var user models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// validate the struct
	if err := validate.Struct(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// check if user with email exists
	emailExists, emailErr := services.UserWithEmailExists(user.Email)
	if emailErr != nil {
		log.Panic(emailErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": emailErr.Error()})
	} else if emailExists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User with this email or phone already exists"})
		return
	}

	// check if user with phone exists
	phoneExists, phoneErr := services.UserWithPhoneExists(user.Phone)
	if phoneErr != nil {
		log.Panic(phoneErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": phoneErr.Error()})
	} else if phoneExists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User with this email or phone already exists"})
		return
	}

	// hash the password
	hashPassword := helpers.Hashpassword(user.Password)
	createdAt := time.Now().UTC().Format(time.RFC3339)

	user.Password = hashPassword
	user.CreatedAt = createdAt
	user.UpdatedAt = createdAt
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()

	// get access and refresh token
	accessToken, refreshToken, err := helpers.GenerateTokens(user);
	if err != nil {
		log.Panic("Token generation failure")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	// create user account
	if err := services.CreateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! user account not created"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": user, "statusCode": http.StatusCreated})
}

func Login(ctx *gin.Context){
	var user models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the data of the user with the email provided
	newUser, err := services.GetUserWithEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Incorrect email or password"})
		return
	}

	// Verify if the user password in database is the same as the password
	// supplied by the user
	passwordIsValid := helpers.VerifyPassword(newUser.Password, user.Password)
	if !passwordIsValid {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Incorrect email or password"})
		return
	}

	// Generate and Update the user tokens in database
	err1 := services.GenerateAndUpdateUserTokens(&newUser);
	if err1 != nil {
		log.Panic("Token generation failure")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": newUser, "statusCode": http.StatusOK})
}

func GetUsers(ctx *gin.Context){
	if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 20
	}

	page, err1 := strconv.Atoi(ctx.Query("page"))
	if err1 != nil || page < 1 {
		page = 1
	}

	// get user with userId
	users, err2 := services.GetUsers(recordPerPage, page)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users, "nextPage": page + 1,  "statusCode": http.StatusOK})

}

func GetUser(ctx *gin.Context){
	userId := ctx.Param("userId")

	if err := helpers.MatchUserTypeToUserID(ctx, userId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error":  err.Error()})
		return
	}

	// get user with userId
	user, err := services.GetUser(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": user, "statusCode": http.StatusOK})
}
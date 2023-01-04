package services

import (
	"context"
	"time"

	"github.com/Aanu1995/golang-authentication/database"
	"github.com/Aanu1995/golang-authentication/helpers"
	"github.com/Aanu1995/golang-authentication/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection("Users")

func UserWithEmailExists(email string) (userExists bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return
	}

	userExists = count > 0
	return
}

func UserWithPhoneExists(phone string) (userExists bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"phone": phone})
	if err != nil {
		return
	}

	userExists = count > 0
	return
}

func GetUser(userId string) (user models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	result := userCollection.FindOne(ctx, bson.M{"userid": userId})

	err = result.Decode(&user)

	return
}

func GetUsers(recordPerPage int, page int) (users []models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * recordPerPage))
	opts.SetLimit(int64(recordPerPage))

	result, err := userCollection.Find(ctx, bson.D{}, opts)
	defer result.Close(context.Background())

	if err != nil {
		return
	}

	if err = result.All(context.Background(), &users); err != nil {
		return
	}

	if users == nil {
		users = []models.User{}
	}
	return
}

func GetUserWithEmail(email string) (user models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
	defer cancel()

	result := userCollection.FindOne(ctx, bson.M{"email": email})

	err = result.Decode(&user)

	return
}

func CreateUser(user models.User) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	_, err = userCollection.InsertOne(ctx, user)

	return
}

func GenerateAndUpdateUserTokens(user *models.User) (err error) {
	accessToken, refreshToken, err := helpers.GenerateTokens(*user);
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	updatedAt := time.Now().UTC().Format(time.RFC3339)
	update := bson.M{"$set": bson.M{"accesstoken": accessToken, "refreshtoken": refreshToken, "updatedat":updatedAt}}
	_, err = userCollection.UpdateOne(ctx, bson.M{"userid": user.UserId}, update)
	if err != nil {
		return
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	user.UpdatedAt = updatedAt

	return
}

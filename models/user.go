package models

import "go.mongodb.org/mongo-driver/bson/primitive"


type User struct {
	ID 							primitive.ObjectID  `bson:"_id"`
	UserId 					string							`json:"userId"`
	FirstName 			string							`json:"firstName" validate:"required,min=2,max=100"`
	LastName 				string							`json:"lastName" validate:"required,min=2,max=100"`
	Email 					string							`json:"email" validate:"required,email"`
	Password				string							`json:"password" validate:"required,min=6"`
	Phone 					string							`json:"phone" validate:"required"`
	UserType 				string							`json:"userType" validate:"required,eq=ADMIN|eq=USER"`
	AccessToken 		string							`json:"accessToken"`
	RefreshToken 		string							`json:"refreshToken"`
	CreatedAt 			string							`json:"createdAt"`
	UpdatedAt 			string							`json:"updatedAt"`
}
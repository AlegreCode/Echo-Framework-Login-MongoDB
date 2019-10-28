package models

import (
	"context"
	"errors"
	"log"

	"github.com/gookit/validate"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	. "github.com/alegrecode/echo/LoginMongoDB/db"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" form:"id,omitempty"`
	Name     string             `json:"name" bson:"name" form:"name" validate:"required"`
	Lastname string             `json:"lastname" bson:"lastname" form:"lastname" validate:"required"`
	Email    string             `json:"email" bson:"email" form:"email" validate:"required|email"`
	Age      string             `json:"age" bson:"age" form:"age" validate:"required|min:1|max:150"`
	Password string             `json:"password" bson:"password" form:"password" validate:"required|minLen:4"`
}

func (f User) ConfigValidation(v *validate.Validation) {
	v.StopOnError = false
	v.WithScenes(validate.SValues{
		"login": []string{"Email", "Password"},
		"register": []string{"Name", "Lastname", "Email", "Age", "Password"},
	})
}

func (f User) Messages() map[string]string {
	return validate.MS{
		"required":    "The {field} field is required.",
		"Email.email": "Email not valid.",
	}
}

func SaveUser(c echo.Context) *mongo.InsertOneResult {
	user := new(User)
	if err := c.Bind(user); err != nil {
		log.Fatal(err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hash)
	// collection := Client.Database("go_auth").Collection("users")
	collection := DB.Collection("users")
	insertResult, err2 := collection.InsertOne(context.TODO(), user)
	if err2 != nil {
		log.Fatal(err2)
	}
	return insertResult
}

func GetSingleUser(email string) (User, error) {
	var user User
	// collection := Client.Database("go_auth").Collection("users")
	collection := DB.Collection("users")
	filter := bson.M{"email": email}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return user, errors.New("User not found.")
	}
	return user, nil
}

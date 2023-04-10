package controllers

import (
	"context"
	"net/http"
	"time"
	"workspace/configs"
	"workspace/models"
	"workspace/responses"
	"workspace/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "Users")
var validate = validator.New()

func Register(c*fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":err.Error()}})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data": validationErr.Error()}})
	}

	hash, _ := utils.CreateHash(user.Password)

	var usr models.User

	mongoErr := UserCollection.FindOne(ctx, bson.M{"username":user.Username}).Decode(&usr)
	if mongoErr == nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":"user alredy exsist"}})
	}

	newUser := models.User {
		Id: primitive.NewObjectID(),
		Username: user.Username,
		Password: string(hash),
		IsAdmin: false,
		IsTerminal: false,
	}

	result, err := UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":err.Error()}})
	} else {
		return c.Status(http.StatusAccepted).JSON(responses.Response{
			Status: http.StatusAccepted,
			Message: "success",
			Data: &fiber.Map{"data": result}})
	}
}

func GetUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	auth, _ := JwtFromHeader(c, fiber.HeaderAuthorization)
	var users []models.User
	defer cancel()

	claims, err := utils.EncodeAccsesToken(auth); if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.Response{
			Status: http.StatusUnauthorized,
			Message: "error",
			Data:  &fiber.Map{"data":err.Error()}})
	}

	objId, _ := primitive.ObjectIDFromHex(claims.Id)

	var user models.User

	usrErr := UserCollection.FindOne(ctx, bson.M{"_id":objId}).Decode(&user); if usrErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":usrErr.Error()}})
	}

	if !user.IsAdmin {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":"permission denied"}})
	}

	result, mongoErr := UserCollection.Find(ctx, bson.M{})
	if mongoErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":mongoErr.Error()}})
	} else {
		for result.Next(ctx) {
			var singleUser models.User
			if err = result.Decode(&singleUser); err != nil {
				return c.Status(http.StatusInternalServerError).JSON(responses.Response{
					Status: http.StatusInternalServerError,
					Message: "error",
					Data: &fiber.Map{"data":err.Error()}})
			}

			users = append(users, singleUser)
		}
		return c.Status(http.StatusAccepted).JSON(responses.Response{
			Status: http.StatusAccepted,
			Message: "success",
			Data: &fiber.Map{"data":users}})
	}

}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	auth, _ := JwtFromHeader(c, fiber.HeaderAuthorization)
	userid := c.Params("userId")
	defer cancel()

	claims, err := utils.EncodeAccsesToken(auth); if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.Response{
			Status: http.StatusUnauthorized,
			Message: "error",
			Data: &fiber.Map{"data": err.Error()}})
	}

	objId, _ := primitive.ObjectIDFromHex(claims.Id)

	var user models.User

	usrErr := UserCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user); if usrErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":  usrErr.Error()}})
	}

	if !user.IsAdmin {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":"permission denied"}})
	}

	objId2, _ := primitive.ObjectIDFromHex(userid)

	result, mongoErr := UserCollection.DeleteOne(ctx, bson.M{"_id":objId2})
	if mongoErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":mongoErr.Error()}})
	} else {
		return c.Status(http.StatusAccepted).JSON(responses.Response{
			Status: http.StatusAccepted,
			Message: "error",
			Data: &fiber.Map{"data":result.DeletedCount}})
	}
}

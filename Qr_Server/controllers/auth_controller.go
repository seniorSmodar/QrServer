package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"workspace/middleware"
	"workspace/models"
	"workspace/responses"
	"workspace/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user)
	err != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest, 
			Message: "error", 
			Data: &fiber.Map{"data": err.Error()}})
    }

    if validationErr := validate.Struct(&user)
	validationErr != nil {
        return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest, 
			Message: "error", 
			Data: &fiber.Map{"data": validationErr.Error()}})
    }

	var userbd models.User

	err := UserCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&userbd) 
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data": err.Error()}})
	}

	match, err := utils.VerifyHash(user.Password, string(userbd.Password))
	if (!match) {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data": "ErrEmptyPassword"}})
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{
			Status: http.StatusInternalServerError, 
			Message: "error", 
			Data: &fiber.Map{"data": err.Error()}})
	} else {
		token, _ := utils.GenerateNewAccessToken(userbd.Id.Hex())
		return c.Status(http.StatusOK).JSON(responses.TokenResponse{
			Access_token: token,  
			AuthScheme: middleware.JWTAuthScheme})
	}
}

func JwtFromHeader(c *fiber.Ctx, header string) (string, error) {
	auth := c.Get(header)
	l := len(middleware.JWTAuthScheme)
	if len(auth) > l+1 && strings.EqualFold(auth[:l], middleware.JWTAuthScheme) {
		return auth[l+1:], nil
	}
	return "", utils.ErrHeader
}
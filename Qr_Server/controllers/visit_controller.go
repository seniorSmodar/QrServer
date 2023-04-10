package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"workspace/configs"
	"workspace/models"
	"workspace/responses"
	"workspace/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var VisitCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "Visits")

func CreateVisit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	auth, _ := JwtFromHeader(c, fiber.HeaderAuthorization)
	code := c.Params("code")
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

	codeterm, redisErr := configs.RedisGet(); if redisErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data": redisErr.Error()}})
	} 
	
	if codeterm != code {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data": "invalid code"},
		})
	}

	newVisit := models.Visit {
		Id: primitive.NewObjectID(),
		Username: user.Username,
		VisitDate: time.Now(),
	}

	upperLimit := time.Now().Add(time.Hour * 12)

    filter := bson.M{
        "username": newVisit.Username,
        "visitdate": bson.M{
            "$lte": upperLimit,
        },
    }

	visit := VisitCollection.FindOne(ctx, filter); if visit != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":"visit alredy exsist"}})
	} else {
		result, err := VisitCollection.InsertOne(ctx, newVisit)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.Response{
				Status: http.StatusBadRequest,
				Message: "error",
				Data: &fiber.Map{"data":err.Error()}})
		} else {
			return c.Status(http.StatusAccepted).JSON(responses.Response{
				Status: http.StatusAccepted,
				Message: "success",
				Data: &fiber.Map{"data":result}})
		}
	}
}

func CreateQr(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	auth, _ := JwtFromHeader(c, fiber.HeaderAuthorization)
	defer cancel()

	claims, err := utils.EncodeAccsesToken(auth); if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.Response{
			Status: http.StatusUnauthorized,
			Message: "error",
			Data: &fiber.Map{"data":err.Error()}})
	}

	objId, _ := primitive.ObjectIDFromHex(claims.Id)

	var user models.User

	usrErr := UserCollection.FindOne(ctx, bson.M{"_id":objId}).Decode(&user); if usrErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data": usrErr.Error()}})
	}

	if !user.IsTerminal {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":"permission denied"}})
	}

	rInt := utils.RandomInt(100000, 999999)
	configs.RedisSet(strconv.Itoa(rInt))
	
	return c.Status(http.StatusAccepted).JSON(responses.Response{
		Status: http.StatusAccepted,
		Message: "succes",
		Data: &fiber.Map{"data":rInt},
	})
}

func GetVisits(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	auth, _ := JwtFromHeader(c, fiber.HeaderAuthorization)
	var visits []models.Visit
	defer cancel()

	claims, err := utils.EncodeAccsesToken(auth); if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.Response{
			Status: http.StatusUnauthorized,
			Message: "error",
			Data: &fiber.Map{"data": err.Error()}})
	}

	objId, _ := primitive.ObjectIDFromHex(claims.Id)

	var user models.User

	usrErr := UserCollection.FindOne(ctx, bson.M{"_id":objId}).Decode(&user); if usrErr != nil {
		return c.Status(http.StatusBadRequest).JSONP(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":usrErr.Error()}})
	}

	if !user.IsAdmin || !user.IsTerminal {
		result, mongoErr := VisitCollection.Find(ctx, bson.M{"username":user.Username})
		if mongoErr != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.Response{
				Status: http.StatusBadRequest,
				Message: "error",
				Data: &fiber.Map{"data":mongoErr.Error()}})
		} else {
			for result.Next(ctx) {
				var singleVisit models.Visit
				if err = result.Decode(&singleVisit); err != nil {
					return c.Status(http.StatusInternalServerError).JSON(responses.Response{
						Status: http.StatusInternalServerError,
						Message: "error",
						Data: &fiber.Map{"data":err.Error()}})
				}
				visits = append(visits, singleVisit)
			}		
			return c.Status(http.StatusAccepted).JSON(responses.Response{
				Status: http.StatusAccepted,
				Message: "success",
				Data: &fiber.Map{"data": visits}})
		}
	} else {
		result, mongoErr := VisitCollection.Find(ctx, bson.M{})
		if mongoErr != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.Response{
				Status: http.StatusBadRequest,
				Message: "error",
				Data: &fiber.Map{"data": mongoErr.Error()}})
		} else {
			for result.Next(ctx) {
				var singleVisit models.Visit
				if err = result.Decode(&singleVisit); err != nil {
					return c.Status(http.StatusInternalServerError).JSON(responses.Response{
						Status: http.StatusInternalServerError,
						Message: "error",
						Data: &fiber.Map{"data": err.Error()}})
				}
				visits = append(visits, singleVisit)
			}
			return c.Status(http.StatusAccepted).JSON(responses.Response{
				Status: http.StatusAccepted,
				Message: "success",
				Data: &fiber.Map{"data": visits}})
		}
	}
}

func DeleteLegacyVisits(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	auth, _ := JwtFromHeader(c, fiber.HeaderAuthorization)
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
			Data: &fiber.Map{"data":usrErr.Error()}})
	}

	if !user.IsAdmin {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data":"permission denied"}})
	}

	filter := bson.M{"visit_date": bson.M{"$lt": time.Now()}}

	result, mongoErr := VisitCollection.DeleteMany(ctx, filter); if mongoErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{
			Status: http.StatusBadRequest,
			Message: "error",
			Data: &fiber.Map{"data": mongoErr.Error()}})
	} else {
		return c.Status(http.StatusAccepted).JSON(responses.Response{
			Status: http.StatusAccepted,
			Message: "success",
			Data: &fiber.Map{"data":result.DeletedCount}})
	}
}
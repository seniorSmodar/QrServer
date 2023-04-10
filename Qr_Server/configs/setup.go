package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: EnvRedisURI(),
		Password: "",
		DB: 0,
	})

	return client
}

var RedisClient *redis.Client = ConnectRedis()

func ConnectDB() *mongo.Client  {
    client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
    if err != nil {
        log.Fatal(err)
    }
  
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Все заебись")
    return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    collection := client.Database("QrServer").Collection(collectionName)
    return collection
}

func RedisSet(message string) {
	RedisClient.Set("pivo", message, 0)
}

func RedisGet() (string, error) {
	value, err := RedisClient.Get("pivo").Result()
	if err != nil{
		return "", err
	}
	return value, nil
}
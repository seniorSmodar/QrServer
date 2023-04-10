package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Введите url базы")
	}

	return os.Getenv("MONGOURI")
}

func EnvPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Введите порт сервера")
	}

	return os.Getenv("PORT")
}

func EnvRedisURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Введите url memory-базы")
	}

	return os.Getenv("REDISURI")
}

func EnvDuration() int {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Введите duration")
	}
	value := os.Getenv("QRDURATION")
	i, err := strconv.Atoi(value)
	if err != nil{
		return 0
	}
	return i
}
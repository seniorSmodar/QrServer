package utils

import (
	"os"
	"strconv"
	"time"
	"errors"
	"math/rand"

	"github.com/golang-jwt/jwt"
	"github.com/matthewhartstonge/argon2"
)

var (
	ErrUnexpected = errors.New("unexpected signing method")
	ErrHeader = errors.New("missing or malformed JWT")
)

type JWTClaims struct {
	jwt.StandardClaims
	Field string `json:"field"`
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}


func CreateHash(password string) ([]byte, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(password))

	if err != nil {
		return nil, err
	}
	
	return encoded, nil
}

func VerifyHash(password, hash string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(hash))
	if err != nil {
		return false, err
	}

	return ok, nil
}

func GenerateNewAccessToken(user_id string) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")

	minutesCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))

	claims := jwt.StandardClaims{
		Id:        user_id,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func EncodeAccsesToken(token string) (*JWTClaims, error) {
	token2, _ := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpected
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	JWTClaims := token2.Claims.(*JWTClaims)

	return JWTClaims, nil
}
package responses

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    *fiber.Map `json:"detail"`
}

type TokenResponse struct {
	Access_token  string `json:"access_token"`
	AuthScheme    string `json:"auth_scheme"`
}
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("RAHASIA_DAPUR_ARTELA") // Pindahkan ke .env nanti

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil Token dari Header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"message": "Unauthorized"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 2. Parse Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"message": "Invalid Token"})
		}

		// Lanjut ke handler berikutnya
		return c.Next()
	}
}
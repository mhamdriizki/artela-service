package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

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
			// FIX: Gunakan logic yang sama dengan AuthHandler
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "RAHASIA_DAPUR_ARTELA" // Default fallback jika env tidak ada
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"message": "Invalid Token"})
		}

		// Lanjut ke handler berikutnya
		return c.Next()
	}
}
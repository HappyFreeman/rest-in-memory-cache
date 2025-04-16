package middleware

import (
	"github.com/HappyFreeman/rest-in-memory-cache/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

// Обычный миддлваер

func Authorization(jwtClient jwt.JWTClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		const bearerPrefix = "Bearer "
		if len(tokenString) <= len(bearerPrefix) || tokenString[:len(bearerPrefix)] != bearerPrefix {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		tokenString = tokenString[len(bearerPrefix):]

		TokenData, err := jwtClient.GetDataFromToken(&jwt.GetDataFromTokenParams{
			Token: tokenString,
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("userId", TokenData.UserId)

		return c.Next()
	}
}

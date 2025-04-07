package middleware

import (
	"github.com/HappyFreeman/rest-in-memory-cache/internal/config"
	"github.com/HappyFreeman/rest-in-memory-cache/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

// Обычный миддлваер

func Authorization(cfgJWT config.JWT) fiber.Handler {
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

		userId, err := jwt.GetUserId(tokenString, cfgJWT.Secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("user_id", userId)

		return c.Next()
	}
}

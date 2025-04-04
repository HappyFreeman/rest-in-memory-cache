package middleware

import (
	"github.com/HappyFreeman/rest-in-memory-cache/internal/config"
	"github.com/HappyFreeman/rest-in-memory-cache/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

// Обычный миддлваер

func Authorization(token string, cfgJWT config.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// проверка токена вторизации
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		// Убираем Bearer перед токеном
		tokenString = tokenString[len("Bearer "):]

		userId, err := jwt.GetUserId(tokenString, cfgJWT.Secret)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Сохраняем user_id в `ctx.Locals`
		c.Locals("user_id", userId)

		return c.Next()
	}
}

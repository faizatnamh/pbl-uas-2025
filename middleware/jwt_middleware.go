package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware validates JWT and stores claims in c.Locals("user_claims")
func JWTMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Missing Authorization header",
		})
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Invalid Authorization format",
		})
	}

	tokenString := parts[1]
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret123"
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"code":        401,
			"message":     "Unauthorized",
			"description": "Invalid or expired token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	// simpan claims ke locals supaya middleware RBAC & handler bisa pakai
	c.Locals("user_claims", claims)

	return c.Next()
}

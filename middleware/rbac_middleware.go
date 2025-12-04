package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"pbluas/app/repository"
)

func RBACMiddleware(c *fiber.Ctx, permRepo *repository.PermissionRepository, requiredPerms ...string) error {
	// Ambil token dari header Authorization
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(401).JSON(fiber.Map{"message": "missing token"})
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"message": "invalid token format"})
	}

	tokenStr := parts[1]

	// Parse JWT
	secret := []byte("secret123")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "invalid or expired token"})
	}

	claims := token.Claims.(jwt.MapClaims)

	role := claims["role"].(string)

	// Admin = full access
	if role == "Admin" {
		return c.Next()
	}

	// Query permission berdasarkan role
	permissions, err := permRepo.GetPermissionsByRole(role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to load permissions"})
	}

	// Cek apakah role memiliki salah satu permission
	for _, needed := range requiredPerms {
		for _, owned := range permissions {
			if owned == needed {
				return c.Next()
			}
		}
	}

	return c.Status(403).JSON(fiber.Map{"message": "forbidden: permission denied"})
}

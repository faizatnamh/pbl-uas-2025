package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"pbluas/app/service"
)

func AchievementRoute(api fiber.Router, achievementService *service.AchievementService) {

	ach := api.Group("/achievements")

	ach.Post("/", func(c *fiber.Ctx) error {

		// ✅ ambil claims dengan tipe YANG BENAR
		claims, ok := c.Locals("user_claims").(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"message": "unauthorized",
			})
		}

		// ambil role
		roleVal, ok := claims["role"]
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"message": "role not found",
			})
		}

		role, ok := roleVal.(string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"message": "invalid role",
			})
		}

		// hanya admin & mahasiswa
		if role != "Admin" && role != "Mahasiswa" {
			return c.Status(403).JSON(fiber.Map{
				"message": "only admin and mahasiswa can create achievement",
			})
		}

		return achievementService.CreateHandler(c)
	})
}

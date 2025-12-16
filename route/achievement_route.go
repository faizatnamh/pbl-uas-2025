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

	ach.Get("/", achievementService.ListByRole)
	ach.Get("/:id", achievementService.Detail)
	ach.Put("/:id", achievementService.Update)
	ach.Post("/:id/attachments", achievementService.UploadAttachment)
	ach.Delete("/:id", achievementService.Delete)
	ach.Post("/:id/submit", achievementService.Submit)
	ach.Post("/:id/verify", achievementService.Verify)
	ach.Post("/:id/reject", achievementService.Reject)
	ach.Get("/:id/history", achievementService.History)
}	

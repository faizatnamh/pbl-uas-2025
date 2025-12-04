package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/middleware"
	"pbluas/app/repository"
)

func MahasiswaRoute(api fiber.Router, permRepo *repository.PermissionRepository) {
	mhs := api.Group("/mahasiswa")

	mhs.Use(middleware.JWTMiddleware)
	mhs.Use(func(c *fiber.Ctx) error {
		return middleware.RBACMiddleware(c, permRepo,
			"achievement:create",
			"achievement:read",
			"achievement:update",
		)
	})

	mhs.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Mahasiswa Only"})
	})
}


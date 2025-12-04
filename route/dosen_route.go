package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/middleware"
	"pbluas/app/repository"
)

func DosenRoute(api fiber.Router, permRepo *repository.PermissionRepository) {
	dosen := api.Group("/dosen")

	dosen.Use(middleware.JWTMiddleware)
	dosen.Use(func(c *fiber.Ctx) error {
		return middleware.RBACMiddleware(c, permRepo,
			"achievement:read",
			"achievement:verify",
		)
	})

	dosen.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Dosen Wali Only"})
	})
}


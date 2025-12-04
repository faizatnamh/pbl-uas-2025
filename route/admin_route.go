package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/middleware"
	"pbluas/app/repository"
)

func AdminRoute(api fiber.Router, permRepo *repository.PermissionRepository) {
	admin := api.Group("/admin")

	admin.Use(middleware.JWTMiddleware)
	admin.Use(func(c *fiber.Ctx) error {
		return middleware.RBACMiddleware(c, permRepo, "*")
	})

	admin.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Admin Only"})
	})
}

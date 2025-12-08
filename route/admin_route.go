package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/app/service"
	"pbluas/middleware"
	"pbluas/app/repository"
)

func AdminRoute(api fiber.Router, permRepo *repository.PermissionRepository, userService *service.UserService) {

	require := func(perms ...string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			return middleware.RBACMiddleware(c, permRepo, perms...)
		}
	}

	// ========== ADMIN MANAGE USERS ==========
	api.Get("/users", require("user:manage"), userService.GetAllUsers)
	api.Get("/users/:id", require("user:manage"), userService.GetUserByID)
	api.Post("/users", require("user:manage"), userService.CreateUser)
	api.Put("/users/:id", require("user:manage"), userService.UpdateUser)
	api.Delete("/users/:id", require("user:manage"), userService.DeleteUser)
	api.Put("/users/:id/role", require("user:manage"), userService.UpdateUserRole)
}

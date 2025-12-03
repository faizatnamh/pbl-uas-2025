package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/app/service"
)

func AuthRoute(router fiber.Router, userService *service.UserService) {
	auth := router.Group("/auth")
	auth.Post("/login", userService.Login)
	auth.Get("/profile", userService.Profile)
}

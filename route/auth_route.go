package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/app/service"
)

func AuthRoute(router fiber.Router, userService *service.UserService) {
	router.Post("/login", userService.Login)
	router.Get("/profile", userService.Profile)
}

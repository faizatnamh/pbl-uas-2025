package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/app/service"
)

func AuthRoute(router fiber.Router, userService *service.UserService) {
	router.Post("/login", userService.Login)
	router.Post("/refresh", userService.Refresh)
    router.Post("/logout", userService.Logout)
	router.Get("/profile", userService.Profile)
}

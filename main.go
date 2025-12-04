package main

import (
    "pbluas/config"
    "pbluas/database"
    "pbluas/app/repository"
    "pbluas/app/service"
    "pbluas/middleware"
    "pbluas/route"

    "github.com/gofiber/fiber/v2"
)

func main() {
    // Load .env
    config.LoadEnv()

    // Fiber app
    app := fiber.New()

    // Database
    db := database.ConnectPostgres()

    // Init repo
    userRepo := repository.NewUserRepository(db)
    permRepo := repository.NewPermissionRepository(db)

    // Init service
    userService := service.NewUserService(userRepo, permRepo)
 
    // PUBLIC ROUTES (NO TOKEN)
    auth := app.Group("/api/v1/auth")
    route.AuthRoute(auth, userService)

 
    // PROTECTED ROUTES (JWT)
    api := app.Group("/api/v1")
    api.Use(middleware.JWTMiddleware)

    route.AdminRoute(api, permRepo)
    route.MahasiswaRoute(api, permRepo)
    route.DosenRoute(api, permRepo)

    // Start server
    app.Listen(":8080")
}

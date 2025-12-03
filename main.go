package main

import (
    "pbluas/config"
    "pbluas/database"
    "pbluas/app/repository"
    "pbluas/app/service"
    "pbluas/route"

    "github.com/gofiber/fiber/v2"
)

func main() {
    config.LoadEnv()

    app := fiber.New()

    // Connect DB
    db := database.ConnectPostgres()

    // Init repo & service
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)

    // =============================
    // REGISTER ROUTE DI SINI!
    // =============================

    api := app.Group("/api/v1")
    route.AuthRoute(api, userService)   // <-- ROUTE LOGIN & PROFILE AKTIF

    app.Listen(":8080")
}

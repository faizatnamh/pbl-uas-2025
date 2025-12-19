package main

import (
    "pbluas/config"
    "pbluas/database"

    "pbluas/app/repository"
    "pbluas/app/service"

    "pbluas/middleware"
    "pbluas/route"

    "github.com/gofiber/fiber/v2"

	_ "pbluas/docs" 
	fiberSwagger "github.com/swaggo/fiber-swagger" 
)

// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0
// @description Backend API untuk manajemen user dan prestasi mahasiswa
// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	config.LoadEnv()

	app := fiber.New()

	// ===== SWAGGER ROUTE =====
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// DB
	db := database.ConnectPostgres()
	database.ConnectMongo()
	mongoDB := database.MongoDB()


	// -------- INIT REPOSITORIES --------
	userRepo := repository.NewUserRepository(db)
	permRepo := repository.NewPermissionRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	lecturerRepo := repository.NewLecturerRepository(db)
	achievementRepo := repository.NewAchievementRepository(mongoDB)
	achievementRefRepo := repository.NewAchievementReferenceRepository(db)

	// -------- INIT SERVICES --------
	userService := service.NewUserService(userRepo, permRepo)
	studentService := service.NewStudentService(studentRepo, lecturerRepo,  achievementRepo, achievementRefRepo )
	lecturerService := service.NewLecturerService(lecturerRepo, studentRepo)
	achievementService := service.NewAchievementService(achievementRepo,achievementRefRepo,studentRepo, )
	reportService := service.NewReportService(studentRepo, achievementRefRepo, achievementRepo)

	// -------- PUBLIC ROUTES --------
	auth := app.Group("/api/v1/auth")
	route.AuthRoute(auth, userService)

	// -------- PROTECTED ROUTES --------
	api := app.Group("/api/v1")
	api.Use(middleware.JWTMiddleware)
	route.AdminRoute(api, permRepo, userService, studentService, lecturerService)
	route.MahasiswaRoute(api, studentService)
	route.AchievementRoute(api, achievementService)
	route.ReportRoutes(api, reportService)


	app.Listen(":8080")
}

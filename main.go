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

// func main() {
//     // Load environment variables
//     config.LoadEnv()

//     // Initialize Fiber
//     app := fiber.New()

//     // Connect to PostgreSQL
//     db := database.ConnectPostgres()
//     database.ConnectMongo()


//     // -------- INIT REPOSITORIES --------
//     userRepo := repository.NewUserRepository(db)
//     permRepo := repository.NewPermissionRepository(db)
//     studentRepo := repository.NewStudentRepository(db)
//     lecturerRepo := repository.NewLecturerRepository(db)


//     // -------- INIT SERVICES --------
//     userService := service.NewUserService(userRepo, permRepo)
//     studentService := service.NewStudentService(studentRepo, lecturerRepo)
//     lecturerService := service.NewLecturerService(lecturerRepo, studentRepo)



//     // -------- PUBLIC ROUTES (NO TOKEN) --------
//     auth := app.Group("/api/v1/auth")
//     route.AuthRoute(auth, userService)

//     // -------- PROTECTED ROUTES (JWT) --------
//     api := app.Group("/api/v1")
//     api.Use(middleware.JWTMiddleware)

//     route.AdminRoute(api,permRepo,userService,studentService,lecturerService,)
//     route.MahasiswaRoute(api, studentService)



//     // -------- START SERVER --------
//     app.Listen(":8080")
// }
func main() {
	config.LoadEnv()

	app := fiber.New()

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
	studentService := service.NewStudentService(studentRepo, lecturerRepo)
	lecturerService := service.NewLecturerService(lecturerRepo, studentRepo)

	achievementService := service.NewAchievementService(
		achievementRepo,
		achievementRefRepo,
        studentRepo, // ✅ studentRepo adalah interface
	)

	// -------- PUBLIC ROUTES --------
	auth := app.Group("/api/v1/auth")
	route.AuthRoute(auth, userService)

	// -------- PROTECTED ROUTES --------
	api := app.Group("/api/v1")
	api.Use(middleware.JWTMiddleware)

	route.AdminRoute(api, permRepo, userService, studentService, lecturerService)
	route.MahasiswaRoute(api, studentService)
	route.AchievementRoute(api, achievementService)


	app.Listen(":8080")
}

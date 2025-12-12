package route

// import (
// 	"github.com/gofiber/fiber/v2"
// 	"pbluas/app/service"
// 	"pbluas/middleware"
// 	"pbluas/app/repository"
// )

// func LecturerRoute(api fiber.Router, permRepo *repository.PermissionRepository, svc *service.LecturerService) {
// 	r := api.Group("/lecturers")

// 	// Hanya admin yang boleh mengakses
// 	r.Use(middleware.JWTMiddleware)
// 	r.Use(func(c *fiber.Ctx) error {
// 		return middleware.RBACMiddleware(c, permRepo, "admin:manage")
// 	})

// 	r.Get("/", svc.GetAll)
// 	r.Get("/:id/advisees", svc.GetAdvisees)
// }

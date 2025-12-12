package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/app/service"
	"pbluas/middleware"
	"pbluas/app/repository"
)

func MahasiswaRoute(api fiber.Router, studentService *service.StudentService, permRepo *repository.PermissionRepository) {

	require := func(perms ...string) fiber.Handler {
		return func(c *fiber.Ctx) error {
			return middleware.RBACMiddleware(c, permRepo, perms...)
		}
	}

	// STUDENT ROUTES (akses admin + dosen wali)
	api.Get("/students", require("student:read"), studentService.GetStudents)
	api.Get("/students/:id", require("student:read"), studentService.GetStudentByID)
}

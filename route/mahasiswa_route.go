package route

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/app/service"
)

func MahasiswaRoute(api fiber.Router, studentService *service.StudentService) {

	// STUDENT ROUTES
	api.Get("/students", studentService.GetStudents)
	api.Get("/students/:id", studentService.GetStudentByID)
	api.Get("/students/:id/achievements", studentService.GetStudentAchievements)
}

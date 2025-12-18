package route

import (
	"pbluas/app/service"
	"pbluas/middleware"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(
	api fiber.Router,
	reportService *service.ReportService,
) {
	report := api.Group(
		"/reports",
		middleware.JWTMiddleware, // ✅ LANGSUNG
	)

	report.Get("/student/:id", reportService.GetStudentReport)
	report.Get("/statistics", reportService.GetStatistics)
}

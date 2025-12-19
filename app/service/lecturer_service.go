package service

import (
	"github.com/gofiber/fiber/v2"
	"pbluas/app/repository"
)

type LecturerService struct {
	LecturerRepo repository.LecturerRepository
	StudentRepo  repository.StudentRepository
}

func NewLecturerService(
	lecturerRepo repository.LecturerRepository,
	studentRepo repository.StudentRepository,
) *LecturerService {
	return &LecturerService{
		LecturerRepo: lecturerRepo,
		StudentRepo:  studentRepo,
	}
}

// GetLecturers godoc
// @Summary Get lecturers list
// @Description Get list of lecturers
// @Tags Lecturers
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /lecturers [get]
func (s *LecturerService) GetAll(c *fiber.Ctx) error {
	list, err := s.LecturerRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(list)
}

// GetLecturerAdvisees godoc
// @Summary Get lecturer advisees
// @Description Get students advised by lecturer
// @Tags Lecturers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Success 200 {array} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /lecturers/{id}/advisees [get]
func (s *LecturerService) GetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	list, err := s.StudentRepo.GetStudentsByAdvisor(lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(list)
}

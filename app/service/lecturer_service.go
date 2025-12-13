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

// GET /lecturers (ADMIN)
func (s *LecturerService) GetAll(c *fiber.Ctx) error {
	list, err := s.LecturerRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(list)
}

// GET /lecturers/:id/advisees (ADMIN)
func (s *LecturerService) GetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	list, err := s.StudentRepo.GetStudentsByAdvisor(lecturerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(list)
}

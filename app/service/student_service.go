package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"pbluas/app/repository"
)

type StudentService struct {
	StudentRepo  repository.StudentRepository
	LecturerRepo repository.LecturerRepository
}
type AssignAdvisorRequest struct {
	LecturerID string `json:"lecturer_id" validate:"required"`
}

func NewStudentService(studentRepo repository.StudentRepository, lecturerRepo repository.LecturerRepository) *StudentService {
	return &StudentService{
		StudentRepo:  studentRepo,
		LecturerRepo: lecturerRepo,
	}
}

// GET /students
func (s *StudentService) GetStudents(c *fiber.Ctx) error {
	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userId := claims["id"].(string)

	if role == "Admin" {
		list, err := s.StudentRepo.GetAllStudents()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(list)
	}

	if role == "Dosen Wali" {
		lecturer, err := s.LecturerRepo.GetLecturerByUserID(userId)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "lecturer profile not found"})
		}

		list, err := s.StudentRepo.GetStudentsByAdvisor(lecturer.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}

		return c.JSON(list)
	}

	return c.Status(403).JSON(fiber.Map{"message": "forbidden"})
}

// GET /students/:id
func (s *StudentService) GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userId := claims["id"].(string)

	student, err := s.StudentRepo.GetStudentByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "student not found"})
	}

	// Admin → bebas akses
	if role == "Admin" {
		return c.JSON(student)
	}

	// Dosen Wali → hanya boleh ambil anak bimbinganna
	if role == "Dosen Wali" {
		lecturer, err := s.LecturerRepo.GetLecturerByUserID(userId)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"message": "lecturer not found"})
		}

		if student.AdvisorID != nil && *student.AdvisorID == lecturer.ID {
			return c.JSON(student)
		}

		return c.Status(403).JSON(fiber.Map{"message": "forbidden"})
	}

	return c.Status(403).JSON(fiber.Map{"message": "forbidden"})
}

func (s *StudentService) AssignAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	var req AssignAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	// cek student ada
	_, err := s.StudentRepo.GetStudentByID(studentID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "student not found"})
	}

	// cek lecturer ada
	_, err = s.LecturerRepo.GetLecturerByID(req.LecturerID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "lecturer not found"})
	}

	// update advisor
	if err := s.StudentRepo.UpdateAdvisor(studentID, req.LecturerID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "advisor assigned successfully",
	})
}


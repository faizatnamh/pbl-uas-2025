package service

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"pbluas/app/models"
	"pbluas/app/repository"
)

type StudentService struct {
	StudentRepo        repository.StudentRepository
	LecturerRepo       repository.LecturerRepository
	AchievementRepo    *repository.AchievementRepository
	AchievementRefRepo *repository.AchievementReferenceRepository
}

type AssignAdvisorRequest struct {
	LecturerID string `json:"lecturer_id" validate:"required"`
}

func NewStudentService(
	studentRepo repository.StudentRepository, 
	lecturerRepo repository.LecturerRepository,
	achievementRepo *repository.AchievementRepository,
	achievementRefRepo *repository.AchievementReferenceRepository,
	) *StudentService {
	return &StudentService{
		StudentRepo:  studentRepo,
		LecturerRepo: lecturerRepo,
		AchievementRepo:    achievementRepo,
		AchievementRefRepo: achievementRefRepo,
	}
}

// GetStudents godoc
// @Summary Get students list
// @Description Get list of students based on role
// @Tags Students
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /students [get]
func (s *StudentService) GetStudents(c *fiber.Ctx) error {
	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userId := claims["id"].(string)

	// Admin → semua mahasiswa
	if role == "Admin" {
		list, err := s.StudentRepo.GetAllStudents()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(list)
	}

	// Dosen → semua mahasiswa
	if role == "Dosen Wali" {
		list, err := s.StudentRepo.GetAllStudents()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		return c.JSON(list)
	}

	// Mahasiswa → data diri sendiri
	if role == "Mahasiswa" {
		list, err := s.StudentRepo.GetAllStudents()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}

		for _, student := range list {
			if student.UserID == userId {
				return c.JSON([]interface{}{student})
			}
		}

		return c.Status(404).JSON(fiber.Map{"message": "student profile not found"})
	}

	return c.Status(403).JSON(fiber.Map{"message": "forbidden"})
}

// GetStudentDetail godoc
// @Summary Get student detail
// @Description Get student detail by ID
// @Tags Students
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /students/{id} [get]
func (s *StudentService) GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userId := claims["id"].(string)

	student, err := s.StudentRepo.GetStudentByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "student not found"})
	}

	// Admin → bebas
	if role == "Admin" {
		return c.JSON(student)
	}

	// Dosen → bebas
	if role == "Dosen Wali" {
		return c.JSON(student)
	}

	// Mahasiswa → hanya diri sendiri
	if role == "Mahasiswa" {
		if student.UserID == userId {
			return c.JSON(student)
		}
		return c.Status(403).JSON(fiber.Map{"message": "forbidden"})
	}

	return c.Status(403).JSON(fiber.Map{"message": "forbidden"})
}

// AssignAdvisor godoc
// @Summary Assign academic advisor
// @Description Assign lecturer as advisor to student
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param body body AssignAdvisorRequest true "Advisor payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /students/{id}/advisor [put]
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

// GetStudentAchievements godoc
// @Summary Get student achievements
// @Description Get achievements of a student
// @Tags Students
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {array} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /students/{id}/achievements [get]
func (s *StudentService) GetStudentAchievements(c *fiber.Ctx) error {
	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userID := claims["id"].(string)

	studentID := c.Params("id")
	if studentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "student id required",
		})
	}

	// ================= RBAC =================
	switch role {

	case "Mahasiswa":
		student, err := s.StudentRepo.GetStudentByUserID(userID)
		if err != nil || student.ID != studentID {
			return c.Status(403).JSON(fiber.Map{
				"message": "forbidden",
			})
		}

	case "Dosen", "Dosen Wali", "Lecturer":
		ok, err := s.AchievementRefRepo.IsAdvisorOfStudent(userID, studentID)
		if err != nil || !ok {
			return c.Status(403).JSON(fiber.Map{
				"message": "forbidden",
			})
		}

	case "Admin":
		// full access

	default:
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	// ================= GET REFERENCES =================
	refs, err := s.AchievementRefRepo.GetByStudentID(studentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if len(refs) == 0 {
		return c.JSON([]fiber.Map{})
	}

	// ================= GET MONGO DATA =================
	var mongoIDs []string
	statusMap := make(map[string]models.AchievementReference)

	for _, r := range refs {
		mongoIDs = append(mongoIDs, r.MongoID)
		statusMap[r.MongoID] = r
	}

	achievements, err := s.AchievementRepo.FindByIDs(context.Background(), mongoIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// ================= RESPONSE =================
	var response []fiber.Map
	for _, a := range achievements {
		ref := statusMap[a.ID.Hex()]

		response = append(response, fiber.Map{
			"id":              a.ID.Hex(),
			"achievementType": a.AchievementType,
			"title":           a.Title,
			"description":     a.Description,
			"details":         a.Details,
			"tags":            a.Tags,
			"status":          ref.Status,
			"submittedAt":     ref.SubmittedAt,
			"verifiedAt":      ref.VerifiedAt,
			"verifiedBy":      ref.VerifiedBy,
			"rejectionNote":   ref.RejectionNote,
			"createdAt":       a.CreatedAt,
		})
	}

	return c.JSON(response)
}

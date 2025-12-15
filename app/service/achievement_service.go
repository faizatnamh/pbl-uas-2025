package service

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"pbluas/app/models"
	"pbluas/app/repository"
)

type AchievementService struct {
	AchievementRepo *repository.AchievementRepository
	ReferenceRepo   *repository.AchievementReferenceRepository
	StudentRepo     repository.StudentRepository 
}

func NewAchievementService(
	ar *repository.AchievementRepository,
	rr *repository.AchievementReferenceRepository,
	sr repository.StudentRepository,

	) *AchievementService {
	return &AchievementService{
		AchievementRepo: ar,
		ReferenceRepo:   rr,
		StudentRepo:     sr,
	}
}

// ===== BUSINESS LOGIC (TIDAK DIUBAH) =====
func (s *AchievementService) Create(ctx context.Context, studentID string, req *models.Achievement) error {
	if studentID == "" {
		return errors.New("invalid student id")
	}

	// Simpan ke Mongo
	req.StudentID = studentID
	if err := s.AchievementRepo.Create(ctx, req); err != nil {
		return err
	}

	// Simpan reference ke Postgres
	ref := &models.AchievementReference{
		StudentID: studentID,
		MongoID:   req.ID.Hex(),
	}

	return s.ReferenceRepo.Create(ref)
}

func (s *AchievementService) CreateHandler(c *fiber.Ctx) error {

	claims, ok := c.Locals("user_claims").(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		return c.Status(401).JSON(fiber.Map{
			"message": "invalid user id",
		})
	}

	var reqBody models.AchievementCreateRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	achievement := &models.Achievement{
		AchievementType: reqBody.AchievementType,
		Title:           reqBody.Title,
		Description:     reqBody.Description,
		Details:         reqBody.Details,
		Tags:            reqBody.Tags,
	}

	// 🔥 INI KUNCI UTAMA
	student, err := s.StudentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "student profile not found",
		})
	}

	// 🔥 PAKAI students.id BUKAN userID
	if err := s.Create(context.Background(), student.ID, achievement); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "achievement created successfully",
	})
}



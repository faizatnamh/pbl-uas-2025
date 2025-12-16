package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"time"
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

func (s *AchievementService) ListByRole(c *fiber.Ctx) error {
	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userID := claims["id"].(string)

	var refs []models.AchievementReference
	var err error

	switch role {
case "Mahasiswa":
	student, err := s.StudentRepo.GetStudentByUserID(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "student profile not found",
		})
	}
	refs, err = s.ReferenceRepo.GetByStudentID(student.ID)

case "Admin":
	refs, err = s.ReferenceRepo.GetAll()

// 🔥 INI KUNCI UTAMA UNTUK DOSEN WALI
case "Dosen", "Dosen Wali", "Lecturer":
	refs, err = s.ReferenceRepo.GetByAdvisorUserID(userID)

default:
	return c.Status(403).JSON(fiber.Map{
		"message": "forbidden",
	})
}


	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// ambil mongo ids
	var mongoIDs []string
	statusMap := make(map[string]string)

	for _, r := range refs {
		mongoIDs = append(mongoIDs, r.MongoID)
		statusMap[r.MongoID] = r.Status
	}

	achievements, err := s.AchievementRepo.FindByIDs(context.Background(), mongoIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// gabungkan status
	var response []fiber.Map
	for _, a := range achievements {
		response = append(response, fiber.Map{
			"id":              a.ID.Hex(),
			"studentId":       a.StudentID,
			"achievementType": a.AchievementType,
			"title":           a.Title,
			"description":     a.Description,
			"details":         a.Details,
			"tags":            a.Tags,
			"status":          statusMap[a.ID.Hex()],
			"createdAt":       a.CreatedAt,
		})
	}

	return c.JSON(response)
}

func (s *AchievementService) Detail(c *fiber.Ctx) error {
	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userID := claims["id"].(string)

	achievementID := c.Params("id")
	if achievementID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "achievement id required",
		})
	}

	// 1️⃣ ambil reference dulu (POSTGRES)
	ref, err := s.ReferenceRepo.GetByMongoID(achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "achievement not found",
		})
	}

	// 2️⃣ RBAC CHECK
	switch role {

	case "Mahasiswa":
		student, err := s.StudentRepo.GetStudentByUserID(userID)
		if err != nil || student.ID != ref.StudentID {
			return c.Status(403).JSON(fiber.Map{
				"message": "forbidden",
			})
		}

	case "Dosen", "Dosen Wali", "Lecturer":
		allowed, err := s.ReferenceRepo.IsAdvisorOfStudent(userID, ref.StudentID)
		if err != nil || !allowed {
			return c.Status(403).JSON(fiber.Map{
				"message": "forbidden",
			})
		}

	case "Admin":
		// bebas

	default:
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	// 3️⃣ ambil detail Mongo
	achievement, err := s.AchievementRepo.FindByID(context.Background(), achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "achievement detail not found",
		})
	}

	// 4️⃣ response gabungan
	return c.JSON(fiber.Map{
		"id":              achievement.ID.Hex(),
		"studentId":       achievement.StudentID,
		"achievementType": achievement.AchievementType,
		"title":           achievement.Title,
		"description":     achievement.Description,
		"details":         achievement.Details,
		"tags":            achievement.Tags,
		"status":          ref.Status,
		"submittedAt":     ref.SubmittedAt,
		"verifiedAt":      ref.VerifiedAt,
		"verifiedBy":      ref.VerifiedBy,
		"rejectionNote":   ref.RejectionNote,
		"createdAt":       achievement.CreatedAt,
	})
}

func (s *AchievementService) Update(c *fiber.Ctx) error {

	// 🔐 ambil claims
	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userID := claims["id"].(string)

	// ❌ hanya mahasiswa
	if role != "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only mahasiswa can update achievement",
		})
	}

	achievementID := c.Params("id")

	// 1️⃣ ambil reference
	ref, err := s.ReferenceRepo.GetByMongoID(achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "achievement not found",
		})
	}

	// 2️⃣ cek owner
	student, err := s.StudentRepo.GetStudentByUserID(userID)
	if err != nil || student.ID != ref.StudentID {
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	// 3️⃣ cek status
	if ref.Status != "draft" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only draft achievement can be updated",
		})
	}

	// 4️⃣ parse body
	var req models.AchievementCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	// 5️⃣ update MongoDB
	err = s.AchievementRepo.UpdateByID(
		context.Background(),
		achievementID,
		req,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "achievement updated successfully",
	})
}

func (s *AchievementService) UploadAttachment(c *fiber.Ctx) error {

	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userID := claims["id"].(string)

	if role != "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only mahasiswa can upload attachment",
		})
	}

	achievementID := c.Params("id")

	// 1️⃣ ambil reference
	ref, err := s.ReferenceRepo.GetByMongoID(achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "achievement not found",
		})
	}

	// 2️⃣ cek owner
	student, err := s.StudentRepo.GetStudentByUserID(userID)
	if err != nil || student.ID != ref.StudentID {
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	// 3️⃣ cek status
	if ref.Status != "draft" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only draft achievement can upload attachment",
		})
	}

	// 4️⃣ ambil file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "file is required",
		})
	}

	// 5️⃣ simpan file
	uploadDir := "./uploads/achievements"
	_ = os.MkdirAll(uploadDir, os.ModePerm)

	filename := fmt.Sprintf("%s_%s", achievementID, file.Filename)
	filePath := path.Join(uploadDir, filename)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to save file",
		})
	}

	// 6️⃣ simpan metadata ke Mongo
	attachment := models.AchievementAttachment{
		FileName:   file.Filename,
		FileURL:    "/uploads/achievements/" + filename,
		FileType:   file.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	err = s.AchievementRepo.AddAttachment(
		context.Background(),
		achievementID,
		attachment,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "attachment uploaded successfully",
	})
}

func (s *AchievementService) Delete(c *fiber.Ctx) error {

	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userID := claims["id"].(string)

	// 1️⃣ hanya mahasiswa
	if role != "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only mahasiswa can delete achievement",
		})
	}

	achievementID := c.Params("id")
	if achievementID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "achievement id required",
		})
	}

	// 2️⃣ ambil reference
	ref, err := s.ReferenceRepo.GetByMongoID(achievementID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "achievement not found",
		})
	}

	// 3️⃣ cek owner
	student, err := s.StudentRepo.GetStudentByUserID(userID)
	if err != nil || student.ID != ref.StudentID {
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	// 4️⃣ cek status
	if ref.Status != "draft" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only draft achievement can be deleted",
		})
	}

	// 5️⃣ SOFT DELETE
	if err := s.ReferenceRepo.SoftDelete(ref.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "achievement deleted successfully",
	})
}

func (s *AchievementService) Submit(c *fiber.Ctx) error {

	// ================= AUTH =================
	claims, ok := c.Locals("user_claims").(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	userID := claims["id"].(string)
	role := claims["role"].(string)
	mongoID := c.Params("id")

	// ================= PERMISSION =================
	if role != "Mahasiswa" && role != "Admin" {
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	// ================= GET REFERENCE =================
	// pakai method yang SUDAH ADA di repo kamu
	ref, err := s.ReferenceRepo.GetByMongoID(mongoID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "achievement not found",
		})
	}

	// ================= STATUS CHECK =================
	if ref.Status != "draft" {
		return c.Status(400).JSON(fiber.Map{
			"message": "only draft achievement can be submitted",
		})
	}

	// ================= OWNERSHIP =================
	if role == "Mahasiswa" {
		student, err := s.StudentRepo.GetStudentByUserID(userID)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "student profile not found",
			})
		}

		if ref.StudentID != student.ID {
			return c.Status(403).JSON(fiber.Map{
				"message": "not your achievement",
			})
		}
	}

	// ================= SUBMIT =================
	// pakai method Submit(id) yang kamu tambahkan di repo
	if err := s.ReferenceRepo.Submit(ref.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "achievement submitted for verification",
	})
}

func (s *AchievementService) Verify(c *fiber.Ctx) error {

	// ================= AUTH =================
	claims := c.Locals("user_claims").(jwt.MapClaims)
	role := claims["role"].(string)
	userID := claims["id"].(string)

	mongoID := c.Params("id")

	// ================= PERMISSION =================
	if role != "Admin" && role != "Dosen" && role != "Dosen Wali" && role != "Lecturer" {
		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}

	// ================= GET REFERENCE =================
	ref, err := s.ReferenceRepo.GetByMongoID(mongoID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "achievement not found",
		})
	}

	// ================= STATUS CHECK =================
	if ref.Status != "submitted" {
		return c.Status(400).JSON(fiber.Map{
			"message": "only submitted achievement can be verified",
		})
	}

	// ================= DOSEN WALI CHECK =================
	if role != "Admin" {
		allowed, err := s.ReferenceRepo.IsAdvisorOfStudent(userID, ref.StudentID)
		if err != nil || !allowed {
			return c.Status(403).JSON(fiber.Map{
				"message": "you are not advisor of this student",
			})
		}
	}

	// ================= VERIFY =================
	if err := s.ReferenceRepo.Verify(ref.ID, userID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "achievement verified successfully",
	})
}

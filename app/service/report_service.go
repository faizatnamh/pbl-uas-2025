package service

import (
	"pbluas/app/models"
	"pbluas/app/repository"
	"sort"
	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	StudentRepo     repository.StudentRepository
	RefRepo         *repository.AchievementReferenceRepository
	AchievementRepo *repository.AchievementRepository
}

func NewReportService(
	studentRepo repository.StudentRepository,
	refRepo *repository.AchievementReferenceRepository,
	achievementRepo *repository.AchievementRepository,
) *ReportService {
	return &ReportService{
		StudentRepo:     studentRepo,
		RefRepo:         refRepo,
		AchievementRepo: achievementRepo,
	}
}

// GetStudentReport godoc
// @Summary Get student achievement report
// @Description Get detailed achievement report of a student
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /reports/student/{id} [get]
func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")

	// 1️⃣ Ambil student detail
	student, err := s.StudentRepo.GetStudentByID(studentID)
	if err != nil {
		return fiber.NewError(404, "student not found")
	}

	// 2️⃣ Ambil achievement references
	refs, err := s.RefRepo.GetByStudentIDForReport(studentID)
	if err != nil {
		return fiber.NewError(500, "failed to load references")
	}

	// 3️⃣ Siapkan response
	var achievements []models.StudentAchievementDTO
	totalPoints := 0
	levelCount := map[string]int{}

	for _, ref := range refs {

		// Ambil detail dari Mongo
		ach, err := s.AchievementRepo.FindByID(c.Context(), ref.MongoID)
		if err != nil {
			continue
		}

	level := ach.Details.CompetitionLevel
	if level == "" {
	level = "unknown"
	}

		achievements = append(achievements, models.StudentAchievementDTO{
			ID:     ref.ID,
			Title:  ach.Title,
			Type:   ach.AchievementType,
			Level:  level,
			Points: ach.Points,
			Status: ref.Status,
		})

		if ref.Status == "verified" {
			totalPoints += ach.Points
		}

		levelCount[level]++
	}

	// 4️⃣ Final response
	return c.JSON(fiber.Map{
		"success": true,
		"data": models.ReportStudentResponse{
			Student: models.StudentReportInfo{
				ID:           student.ID,
				StudentID:    student.StudentID,
				Name:         student.FullName,
				ProgramStudy: student.ProgramStudy,
				AcademicYear: student.AcademicYear,
			},
			Achievements: achievements,
			Summary: models.StudentSummary{
				TotalAchievements: len(achievements),
				TotalPoints:       totalPoints,
				CompetitionLevels: levelCount,
			},
		},
	})
}

// GetStatistics godoc
// @Summary Get achievement statistics
// @Description Get global achievement statistics and analytics
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /reports/statistics [get]
func (s *ReportService) GetStatistics(c *fiber.Ctx) error {
	refs, err := s.RefRepo.GetAll()
	if err != nil {
		return fiber.NewError(500, "failed to load achievement references")
	}

	// =========================
	// 2️⃣ Siapkan variabel statistik
	// =========================
	total := 0

	perType := map[string]int{}
	perMonth := map[string]int{}
	competitionLevels := map[string]int{}
	studentPoints := map[string]int{} // studentID -> points

	// =========================
	// 3️⃣ Loop reference
	// =========================
	for _, ref := range refs {

		// HANYA HITUNG YANG VERIFIED
		if ref.Status != "verified" {
			continue
		}

		// Ambil detail achievement dari Mongo
		ach, err := s.AchievementRepo.FindByID(c.Context(), ref.MongoID)
		if err != nil {
			continue
		}

		total++

		// -------------------------
		// TYPE
		// -------------------------
		perType[ach.AchievementType]++

		// -------------------------
		// PER BULAN (YYYY-MM)
		// -------------------------
		month := ach.CreatedAt.Format("2006-01")
		perMonth[month]++

		// -------------------------
		// COMPETITION LEVEL
		// -------------------------
		level := "unknown"
		if ach.Details.CompetitionLevel != "" {
			level = ach.Details.CompetitionLevel
		}
		competitionLevels[level]++

		// -------------------------
		// POINT MAHASISWA
		// -------------------------
		studentPoints[ref.StudentID] += ach.Points
	}

	// =========================
	// 4️⃣ TOP MAHASISWA (SORT)
	// =========================
	type TopStudent struct {
		StudentID string `json:"student_id"`
		Name      string `json:"name"`
		Points    int    `json:"points"`
	}

	var topStudents []TopStudent

	for studentID, points := range studentPoints {

		student, err := s.StudentRepo.GetStudentByID(studentID)
		if err != nil {
			continue
		}

		topStudents = append(topStudents, TopStudent{
			StudentID: studentID,
			Name:      student.FullName,
			Points:    points,
		})
	}

	// SORT DESC
	sort.Slice(topStudents, func(i, j int) bool {
		return topStudents[i].Points > topStudents[j].Points
	})

	// AMBIL MAX 5
	if len(topStudents) > 5 {
		topStudents = topStudents[:5]
	}

	// =========================
	// 5️⃣ RESPONSE
	// =========================
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total":              total,
			"per_type":           perType,
			"per_month":          perMonth,
			"competition_levels": competitionLevels,
			"top_students":       topStudents,
		},
	})
}


package models

type ReportStudentResponse struct {
	Student     StudentReportInfo        `json:"student"`
	Achievements []StudentAchievementDTO `json:"achievements"`
	Summary     StudentSummary           `json:"summary"`
}

type StudentReportInfo struct {
	ID            string `json:"id"`
	StudentID     string `json:"student_id"`
	Name          string `json:"name"`
	ProgramStudy  string `json:"program_study"`
	AcademicYear  string `json:"academic_year"`
}

type StudentAchievementDTO struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Level  string `json:"level"`
	Points int    `json:"points"`
	Status string `json:"status"`
}

type StudentSummary struct {
	TotalAchievements int            `json:"total_achievements"`
	TotalPoints       int            `json:"total_points"`
	CompetitionLevels map[string]int `json:"competition_levels"`
}

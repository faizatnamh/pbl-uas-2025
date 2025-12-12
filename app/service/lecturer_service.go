package service

import (
	"pbluas/app/models"
	"pbluas/app/repository"
)

type LecturerService struct {
	Repo repository.LecturerRepository
}

func NewLecturerService(repo repository.LecturerRepository) *LecturerService {
	return &LecturerService{Repo: repo}
}

func (s *LecturerService) GetByUserID(userID string) (*models.Lecturer, error) {
	return s.Repo.GetLecturerByUserID(userID)
}

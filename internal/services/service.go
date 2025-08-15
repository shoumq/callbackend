package services

import (
	"database/sql"
	"vasek/internal/dto"
	"vasek/internal/repositories"
)

type RequestService struct {
	db   *sql.DB
	repo repositories.RequestRepository
}

func NewRequestService(db *sql.DB, repo repositories.RequestRepository) *RequestService {
	return &RequestService{
		db:   db,
		repo: repo,
	}
}

func (s *RequestService) CreateRequest(name, text, phone string) (int, error) {
	id, err := s.repo.InsertRequest(name, text, phone)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *RequestService) GetRequest() ([]dto.Request, error) {
	requests, err := s.repo.SelectRequests()
	if err != nil {
		return nil, err
	}
	return requests, nil
}

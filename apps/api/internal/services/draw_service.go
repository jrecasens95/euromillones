package services

import (
	"euromillones/internal/models"
	"euromillones/internal/repositories"
	"euromillones/internal/validators"
)

type DrawService struct {
	repo *repositories.DrawRepository
}

func NewDrawService(repo *repositories.DrawRepository) *DrawService {
	return &DrawService{repo: repo}
}

func (s *DrawService) List(page int, limit int) (repositories.PaginatedDraws, error) {
	return s.repo.List(page, limit)
}

func (s *DrawService) All() ([]models.Draw, error) {
	return s.repo.All()
}

func (s *DrawService) Find(id uint) (models.Draw, error) {
	return s.repo.Find(id)
}

func (s *DrawService) Create(input validators.DrawInput) (models.Draw, error) {
	draw, err := validators.NormalizeDraw(input)
	if err != nil {
		return models.Draw{}, err
	}
	return s.repo.Create(draw)
}

func (s *DrawService) Update(id uint, input validators.DrawInput) (models.Draw, error) {
	draw, err := validators.NormalizeDraw(input)
	if err != nil {
		return models.Draw{}, err
	}
	return s.repo.Update(id, draw)
}

func (s *DrawService) Delete(id uint) error {
	return s.repo.Delete(id)
}

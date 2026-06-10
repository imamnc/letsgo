package health

import "github.com/gofiber/fiber/v2"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Check(ctx *fiber.Ctx) bool {
	if err := s.repository.CheckDatabase(ctx); err != nil {
		return false
	}
	return true
}

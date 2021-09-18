package diary

import (
	"github.com/baryon-m/dear-diary/domain/entity"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(e *Entry) (entity.ID, error) {
	e.ID = entity.NewID()
	e.CreatedAt = time.Now().Format(time.RFC3339)

	return s.repo.Create(e)
}

func (s *Service) FetchOne(id entity.ID) (*Entry, error) {
	return s.repo.FetchOne(id)
}

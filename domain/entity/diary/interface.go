package diary

import "github.com/baryon-m/dear-diary/domain/entity"

type Reader interface {
	FetchOne(id entity.ID) (*Entry, error)
}

type Writer interface {
	Create(e *Entry) (entity.ID, error)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Reader
	Writer
}

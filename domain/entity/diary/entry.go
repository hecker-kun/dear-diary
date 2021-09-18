package diary

import "github.com/baryon-m/dear-diary/domain/entity"

type Entry struct {
	ID entity.ID
	Author string
	Content string
	CreatedAt string
}

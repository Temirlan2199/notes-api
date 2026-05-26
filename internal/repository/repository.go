package repository

import "notes-api/internal/model"

type NoteRepository interface {
	Create(title, content string, tags []string) *model.Note
	GetAll() []*model.Note
	GetByID(id int64) (*model.Note, error)
	Update(id int64, title, content string, tags []string) (*model.Note, error)
	Delete(id int64) error
	Count() int
}

//var _ NoteRepository = (*MemoryRepository)(nil)

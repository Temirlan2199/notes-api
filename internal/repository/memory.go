package repository

import (
	"errors"
	"sort"
	"sync"
	"time"

	"notes-api/internal/model"
)

var ErrNotFound = errors.New("note not found")

type MemoryRepository struct {
	mu     sync.RWMutex
	notes  map[int64]*model.Note
	nextID int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		notes:  make(map[int64]*model.Note),
		nextID: 1,
	}
}

func (r *MemoryRepository) Create(title, content string, tags []string) *model.Note {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	note := &model.Note{
		ID:        r.nextID,
		Title:     title,
		Content:   content,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.notes[r.nextID] = note
	r.nextID++
	return note
}

func (r *MemoryRepository) GetAll() []*model.Note {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Note, 0, len(r.notes))
	for _, note := range r.notes {
		result = append(result, note)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func (r *MemoryRepository) GetByID(id int64) (*model.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	note, ok := r.notes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return note, nil
}

func (r *MemoryRepository) Update(id int64, title, content string, tags []string) (*model.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	note, ok := r.notes[id]
	if !ok {
		return nil, ErrNotFound
	}
	note.Title = title
	note.Content = content
	note.Tags = tags
	note.UpdatedAt = time.Now()
	return note, nil
}

func (r *MemoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.notes[id]; !ok {
		return ErrNotFound
	}
	delete(r.notes, id)
	return nil
}

func (r *MemoryRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.notes)
}

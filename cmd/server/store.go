package main

import (
	"errors"
	"sync"
	"time"
)

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrNotFound = errors.New("note not found")

type Store struct {
	mu     sync.RWMutex
	notes  map[int64]*Note
	nextID int64
}

func NewStore() *Store {
	return &Store{
		notes:  make(map[int64]*Note),
		nextID: 1,
	}
}

func (s *Store) Create(title, content string, tags []string) *Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	note := &Note{
		ID:        s.nextID,
		Title:     title,
		Content:   content,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.notes[s.nextID] = note
	s.nextID++
	return note
}

func (s *Store) GetAll() []*Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Note, 0, len(s.notes))
	for _, note := range s.notes {
		result = append(result, note)
	}

	return result
}

func (s *Store) GetByID(id int64) (*Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	note, ok := s.notes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return note, nil
}

func (s *Store) Update(id int64, title string, content string, tags []string) (*Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	note, ok := s.notes[id]
	if !ok {
		return nil, ErrNotFound
	}
	note.Title = title
	note.Content = content
	note.Tags = tags
	note.UpdatedAt = time.Now()
	return note, nil
}

func (s *Store) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.notes[id]; !ok {
		return ErrNotFound
	}
	delete(s.notes, id)
	return nil
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.notes)
}

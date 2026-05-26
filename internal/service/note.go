package service

import (
	"errors"
	"strings"

	"notes-api/internal/model"
	"notes-api/internal/repository"
)

var (
	ErrTitleRequired  = errors.New("title is required")
	ErrTitleTooLong   = errors.New("title must be at most 200 characters")
	ErrContentTooLong = errors.New("content must be at most 10000 characters")
)

const (
	maxTitleLen   = 200
	maxContentLen = 10000
)

type NoteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(title, content string, tags []string) (*model.Note, error) {
	title = strings.TrimSpace(title)

	if err := validateTitle(title); err != nil {
		return nil, err
	}
	if err := validateContent(content); err != nil {
		return nil, err
	}

	tags = cleanTags(tags)

	return s.repo.Create(title, content, tags), nil
}

func (s *NoteService) List() []*model.Note {
	return s.repo.GetAll()
}

func (s *NoteService) GetByID(id int64) (*model.Note, error) {
	return s.repo.GetByID(id)
}

func (s *NoteService) Update(id int64, title, content string, tags []string) (*model.Note, error) {
	title = strings.TrimSpace(title)

	if err := validateTitle(title); err != nil {
		return nil, err
	}
	if err := validateContent(content); err != nil {
		return nil, err
	}

	tags = cleanTags(tags)

	return s.repo.Update(id, title, content, tags)
}

func (s *NoteService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *NoteService) Count() int {
	return s.repo.Count()
}

func validateTitle(title string) error {
	if title == "" {
		return ErrTitleRequired
	}
	if len(title) > maxTitleLen {
		return ErrTitleTooLong
	}
	return nil
}

func validateContent(content string) error {
	if len(content) > maxContentLen {
		return ErrContentTooLong
	}
	return nil
}

func cleanTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	result := make([]string, 0, len(tags))
	seen := make(map[string]bool)
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		result = append(result, t)
	}
	return result
}

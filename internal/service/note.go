package service

import (
	"errors"
	"fmt"

	"github.com/JannisK89/notes-api/internal/models"
	"github.com/JannisK89/notes-api/internal/repository"
)

var (
	ErrInvalidId   = errors.New("Id must be greater than 0")
	ErrInvalidNote = errors.New("Note must have title and content")
)

type ServiceError struct {
	Src string
	Id  int
	Err error
}

func (e *ServiceError) Error() string {
	if e.Id == 0 {
		return fmt.Sprintf("NoteService error in %s: %v", e.Src, e.Err)
	}
	return fmt.Sprintf("NoteService error in %s with id %v: %v", e.Src, e.Id, e.Err)
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *noteService {
	return &noteService{repo}
}

func (s *noteService) Get(id int) (*models.Note, error) {
	if id < 1 {
		return nil, &ServiceError{"GetNote", id, fmt.Errorf("%w: %v", ErrInvalidId, id)}
	}
	return s.repo.Get(id)
}

func (s *noteService) Create(note *models.Note) (int, error) {
	if note == nil || note.Title == "" || note.Content == "" {
		return 0, &ServiceError{Src: "CreateNote", Err: fmt.Errorf("%w: %v", ErrInvalidNote, note)}
	}
	return s.repo.Create(note)
}

func (s *noteService) GetAll() ([]*models.Note, error) {
	return s.repo.GetAll()
}

func (s *noteService) Update(id int, note *models.Note) error {
	if id < 1 {
		return &ServiceError{"GetNote", id, fmt.Errorf("%w: %v", ErrInvalidId, id)}
	}

	if note == nil || note.Title == "" || note.Content == "" {
		return &ServiceError{"CreateNote", 0, ErrInvalidNote}
	}
	return s.repo.Update(id, note)
}

func (s *noteService) Delete(id int) error {
	if id < 1 {
		return &ServiceError{"GetNote", id, fmt.Errorf("%w: %v", ErrInvalidId, id)}
	}
	return s.repo.Delete(id)
}

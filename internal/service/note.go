package service

import (
	"errors"
	"fmt"

	"github.com/JannisK89/notes-api/internal/models"
	"github.com/JannisK89/notes-api/internal/repository"
)

var (
	ErrInvalidID   = errors.New("ID must be greater than 0")
	ErrInvalidNote = errors.New("Note must have title and content")
)

type ServiceError struct {
	Src string
	ID  int
	Err error
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("NoteService error in %s with id %v: %v", e.Src, e.ID, e.Err)
}

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *noteService {
	return &noteService{repo}
}

func (s *noteService) GetNote(id int) (*models.Note, error) {
	if id < 1 {
		return nil, &ServiceError{"GetNote", id, ErrInvalidID}
	}
	return s.repo.GetNoteByID(id)
}

func (s *noteService) CreateNote(note *models.Note) (int, error) {
	if note == nil || note.Title == "" || note.Content == "" {
		return 0, &ServiceError{"CreateNote", 0, ErrInvalidNote}
	}
	return s.repo.CreateNote(note)
}

func (s *noteService) GetAllNotes() ([]*models.Note, error) {
	return s.repo.GetAllNotes()
}

func (s *noteService) UpdateNote(id int, note *models.Note) error {
	if id < 1 {
		return &ServiceError{"GetNote", id, ErrInvalidID}
	}
	return s.repo.UpdateNoteByID(id, note)
}

func (s *noteService) DeleteNote(id int) error {
	if id < 1 {
		return &ServiceError{"GetNote", id, ErrInvalidID}
	}
	return s.repo.DeleteNoteByID(id)
}

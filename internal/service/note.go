package service

import (
	"fmt"

	"github.com/JannisK89/notes-api/internal/models"
	"github.com/JannisK89/notes-api/internal/repository"
)

type noteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *noteService {
	return &noteService{repo}
}

func (s *noteService) GetNote(id int) (*models.Note, error) {
	if id < 1 {
		return nil, fmt.Errorf("Error due to invalid id: %d", id)
	}
	return s.repo.GetNoteByID(id)
}

func (s *noteService) CreateNote(note *models.Note) (int, error) {
	return s.repo.CreateNote(note)
}

func (s *noteService) GetAllNotes() ([]*models.Note, error) {
	return s.repo.GetAllNotes()
}

func (s *noteService) UpdateNote(id int, note *models.Note) error {
	return s.repo.UpdateNoteByID(id, note)
}

func (s *noteService) DeleteNote(id int) error {
	return s.repo.DeleteNoteByID(id)
}

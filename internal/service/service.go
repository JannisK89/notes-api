package service

import "github.com/JannisK89/notes-api/internal/models"

type NoteService interface {
	Get(id int) (*models.Note, error)
	Create(note *models.Note) (int, error)
	GetAll() ([]*models.Note, error)
	Update(id int, note *models.Note) error
	Delete(id int) error
}

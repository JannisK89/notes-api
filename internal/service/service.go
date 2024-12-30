package service

import "github.com/JannisK89/notes-api/internal/models"

type NoteService interface {
	GetNote(id int) (*models.Note, error)
	CreateNote(note *models.Note) (int, error)
	GetAllNotes() ([]*models.Note, error)
	UpdateNote(id int, note *models.Note) error
	DeleteNote(id int) error
}

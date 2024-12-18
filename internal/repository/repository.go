package repository

import "github.com/JannisK89/notes-api/internal/models"

type NoteRepository interface {
	GetNoteByID(id int) (*models.Note, error)
	CreateNote(note *models.Note) (int, error)
	GetAllNotes() ([]*models.Note, error)
	UpdateNoteByID(id int, note *models.Note) error
	DeleteNoteByID(id int) error
}

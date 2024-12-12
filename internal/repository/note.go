package repository

import (
	"database/sql"

	"github.com/JannisK89/notes-api/internal/models"
)

type noteRepository struct {
	db *sql.DB
}

func NewNotesRepository(db *sql.DB) *noteRepository {
	return &noteRepository{db}
}

func (r *noteRepository) CreateNote(note *models.Note) error {
	_, err := r.db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", note.Title, note.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *noteRepository) GetNoteByID(id int) (*models.Note, error) {
	row := r.db.QueryRow("SELECT id, title, content, FROM notes WHERE id = ?", id)
	note := &models.Note{}
	err := row.Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		return nil, err
	}
	return note, nil
}

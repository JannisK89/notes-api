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
	row := r.db.QueryRow("SELECT id, title, content FROM notes WHERE id = ?", id)
	note := &models.Note{}
	err := row.Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (r *noteRepository) GetAllNotes() ([]*models.Note, error) {
	rows, err := r.db.Query("SELECT id, title, content FROM notes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []*models.Note{}
	for rows.Next() {
		note := &models.Note{}
		err := rows.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *noteRepository) UpdateNoteByID(id int, note *models.Note) error {
	_, err := r.db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", note.Title, note.Content, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *noteRepository) DeleteNoteByID(id int) error {
	_, err := r.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

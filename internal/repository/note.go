package repository

import (
	"database/sql"
	"fmt"

	"github.com/JannisK89/notes-api/internal/models"
)

type RepoError struct {
	Src string
	ID  int
	Err error
}

func (e *RepoError) Error() string {
	return fmt.Sprintf("NoteRepo Error in %s with id %d: %v", e.Src, e.ID, e.Err)
}

func (e *RepoError) Unwrap() error {
	return e.Err
}

type noteRepository struct {
	db *sql.DB
}

func NewNotesRepository(db *sql.DB) *noteRepository {
	return &noteRepository{db}
}

func (r *noteRepository) GetNoteByID(id int) (*models.Note, error) {
	row := r.db.QueryRow("SELECT id, title, content FROM notes WHERE id = ?", id)
	note := &models.Note{}
	err := row.Scan(&note.ID, &note.Title, &note.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &RepoError{"GetNoteByID", id, fmt.Errorf("Note not found", id)}
		}
		return nil, &RepoError{"GetNoteByID", id, err}
	}
	return note, nil
}

func (r *noteRepository) GetAllNotes() ([]*models.Note, error) {
	rows, err := r.db.Query("SELECT id, title, content FROM notes")
	if err != nil {
		return nil, &RepoError{Src: "GetAllNotes", Err: err}
	}
	defer rows.Close()

	notes := []*models.Note{}
	for rows.Next() {
		note := &models.Note{}
		err := rows.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, &RepoError{"GetAllNotes", 0, fmt.Errorf("No notes found")}
			}
			return nil, &RepoError{"GetAllNotes", 0, err}
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *noteRepository) CreateNote(note *models.Note) (int, error) {
	res, err := r.db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", note.Title, note.Content)
	if err != nil {
		return 0, &RepoError{"CreateNote", 0, err}
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, &RepoError{"CreateNote", 0, err}
	}
	return int(id), nil
}

func (r *noteRepository) UpdateNoteByID(id int, note *models.Note) error {
	_, err := r.db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", note.Title, note.Content, id)
	if err != nil {
		return &RepoError{"UpdateNoteByID", id, err}
	}
	return nil
}

func (r *noteRepository) DeleteNoteByID(id int) error {
	_, err := r.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return &RepoError{"DeleteNoteByID", id, err}
	}
	return nil
}

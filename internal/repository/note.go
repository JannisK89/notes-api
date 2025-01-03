package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/JannisK89/notes-api/internal/models"
)

// ErrNoteNotFound is returned when a note with the given ID is not found in the database.
var ErrNoteNotFound = errors.New("note not found")

// RepoError represents an error that occurred within the repository layer.
// It wraps the underlying error and provides additional context, such as the
// source of the error and the ID of the note involved.
type RepoError struct {
	Src string
	Id  int
	Err error
}

func (e *RepoError) Error() string {
	return fmt.Sprintf("NoteRepo error in %s with id %d: %v", e.Src, e.Id, e.Err)
}

func (e *RepoError) Unwrap() error {
	return e.Err
}

// noteRepository implements the NoteRepository interface.
type noteRepository struct {
	db *sql.DB
}

// NewNotesRepository creates a new noteRepository.
func NewNotesRepository(db *sql.DB) *noteRepository {
	return &noteRepository{db}
}

// Get retrieves a note by its ID from the database.
// It returns ErrNoteNotFound if the note is not found.
func (r *noteRepository) Get(id int) (*models.Note, error) {
	row := r.db.QueryRow("SELECT id, title, content FROM notes WHERE id = ?", id)
	note := &models.Note{}
	err := row.Scan(&note.Id, &note.Title, &note.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &RepoError{"GetNoteByID", id, fmt.Errorf("%w: %v", ErrNoteNotFound, err)}
		}
		return nil, &RepoError{"GetNoteByID", id, fmt.Errorf("DB Error: %w", err)}
	}
	return note, nil
}

// GetAll retrieves all notes from the database.
func (r *noteRepository) GetAll() ([]*models.Note, error) {
	rows, err := r.db.Query("SELECT id, title, content FROM notes")
	if err != nil {
		return nil, &RepoError{Src: "GetAllNotes", Err: fmt.Errorf("DB Error: %w", err)}
	}
	defer rows.Close()

	notes := []*models.Note{}
	for rows.Next() {
		note := &models.Note{}
		err := rows.Scan(&note.Id, &note.Title, &note.Content)
		if err != nil {
			return nil, &RepoError{Src: "GetAllNotes", Err: fmt.Errorf("Error Scanning: %w", err)}
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// Create adds a new note to the database.
func (r *noteRepository) Create(note *models.Note) (int, error) {
	res, err := r.db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", note.Title, note.Content)
	if err != nil {
		return 0, &RepoError{Src: "CreateNote", Err: fmt.Errorf("DB Error: %w", err)}
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, &RepoError{Src: "CreateNote", Err: fmt.Errorf("Error getting last Id: %w", err)}
	}
	return int(id), nil
}

// Update modifies an existing note in the database.
func (r *noteRepository) Update(id int, note *models.Note) error {
	_, err := r.db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", note.Title, note.Content, id)
	if err != nil {
		return &RepoError{"UpdateNoteByID", id, fmt.Errorf("DB Error: %w", err)}
	}
	return nil
}

// Delete removes a note from the database.
func (r *noteRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return &RepoError{"DeleteNoteByID", id, fmt.Errorf("DB Error: %w", err)}
	}
	return nil
}

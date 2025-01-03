package service

import (
	"errors"
	"fmt"

	"github.com/JannisK89/notes-api/internal/models"
	"github.com/JannisK89/notes-api/internal/repository"
)

var (
	// ErrInvalidId is returned when a provided ID is invalid.
	ErrInvalidId = errors.New("id must be greater than 0")
	// ErrInvalidNote is returned when the note data is invalid.
	ErrInvalidNote = errors.New("note must have title and content")
)

// Error represents an error that occurred within the service layer. It
// wraps the underlying error and provides additional context, such as the
// source of the error and the ID of the note involved.
type Error struct {
	Src string
	ID  int
	Err error
}

func (e *Error) Error() string {
	if e.ID == 0 {
		return fmt.Sprintf("NoteService error in %s: %v", e.Src, e.Err)
	}
	return fmt.Sprintf("NoteService error in %s with id %v: %v", e.Src, e.ID, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// noteService implements the NoteService interface.
type noteService struct {
	repo repository.NoteRepository
}

// NewNoteService creates a new noteService.
func NewNoteService(repo repository.NoteRepository) *noteService {
	return &noteService{repo}
}

// Get retrieves a note by its ID from the repository.
// It returns ErrInvalidId if the ID is less than 1.
func (s *noteService) Get(id int) (*models.Note, error) {
	if id < 1 {
		return nil, &Error{"GetNote", id, fmt.Errorf("%w: %v", ErrInvalidId, id)}
	}
	return s.repo.Get(id)
}

// Create adds a new note to the repository.
// It returns ErrInvalidNote if the note is nil or if the title or content is
// empty.
func (s *noteService) Create(note *models.Note) (int, error) {
	if note == nil || note.Title == "" || note.Content == "" {
		return 0, &Error{Src: "CreateNote", Err: fmt.Errorf("%w: %v", ErrInvalidNote, note)}
	}
	return s.repo.Create(note)
}

// GetAll retrieves all notes from the repository.
func (s *noteService) GetAll() ([]*models.Note, error) {
	return s.repo.GetAll()
}

// Update modifies an existing note in the repository.
// It returns ErrInvalidId if the ID is less than 1.
// It returns ErrInvalidNote if the note is nil or if the title or content is
// empty.
func (s *noteService) Update(id int, note *models.Note) error {
	if id < 1 {
		return &Error{"GetNote", id, fmt.Errorf("%w: %v", ErrInvalidId, id)}
	}

	if note == nil || note.Title == "" || note.Content == "" {
		return &Error{"CreateNote", 0, ErrInvalidNote}
	}
	return s.repo.Update(id, note)
}

// Delete removes a note from the repository.
// It returns ErrInvalidId if the ID is less than 1.
func (s *noteService) Delete(id int) error {
	if id < 1 {
		return &Error{"GetNote", id, fmt.Errorf("%w: %v", ErrInvalidId, id)}
	}
	return s.repo.Delete(id)
}

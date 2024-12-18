package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/JannisK89/notes-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNoteRepository_CreateNote(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	firstNote := &models.Note{
		Title:   "First Note",
		Content: "This is the first note",
	}

	secondNote := &models.Note{
		Title:   "Second Note",
		Content: "This is the second note",
	}

	repo := NewNotesRepository(db)

	mock.ExpectExec("INSERT INTO notes").WithArgs(firstNote.Title, firstNote.Content).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO notes").WithArgs(secondNote.Title, secondNote.Content).WillReturnResult(sqlmock.NewResult(2, 1))

	// Act
	firstID, firstErr := repo.CreateNote(firstNote)
	secondID, secondErr := repo.CreateNote(secondNote)

	// Assert
	assert.NoError(t, firstErr)
	assert.NoError(t, secondErr)
	assert.Equal(t, 1, firstID)
	assert.Equal(t, 2, secondID)

}

func TestNoteRepository_GetNoteByID(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening a stub database connection: %s", err)
	}
	defer db.Close()
	notes := [2]*models.Note{
		&models.Note{
			ID:      1,
			Title:   "First Note",
			Content: "This is the first note",
		},
		&models.Note{
			ID:      2,
			Title:   "Second Note",
			Content: "This is the second note",
		}}

	repo := NewNotesRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "content"})
	for _, note := range notes {
		rows.AddRow(note.ID, note.Title, note.Content)
	}

	mock.ExpectQuery("SELECT id, title, content FROM notes WHERE id = ?").WithArgs(notes[0].ID).WillReturnRows(rows)
	mock.ExpectQuery("SELECT id, title, content FROM notes WHERE id = ?").WithArgs(notes[1].ID).WillReturnRows(rows)

	// Act
	firstResult, firstErr := repo.GetNoteByID(notes[0].ID)
	secondResult, secondErr := repo.GetNoteByID(notes[1].ID)

	// Assert
	assert.NoError(t, firstErr)
	assert.Equal(t, notes[0], firstResult)
	assert.NoError(t, secondErr)
	assert.Equal(t, notes[1], secondResult)
}

func TestNoteRepository_GetAllNotes(t *testing.T) {

	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening a stub database connection: %s", err)
	}
	defer db.Close()
	notes := [3]*models.Note{
		&models.Note{
			ID:      1,
			Title:   "First Note",
			Content: "This is the first note",
		},
		&models.Note{
			ID:      2,
			Title:   "Second Note",
			Content: "This is the second note",
		},
		&models.Note{
			ID:      3,
			Title:   "Third Note",
			Content: "This is the third note",
		},
	}

	repo := NewNotesRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "content"})
	for _, note := range notes {
		rows.AddRow(note.ID, note.Title, note.Content)
	}

	mock.ExpectQuery("SELECT id, title, content FROM notes").WillReturnRows(rows)

	// Act
	res, err := repo.GetAllNotes()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, len(res), 3)
	for i := range notes {
		assert.Equal(t, res[i], notes[i])
	}
}

func TestNoteRepository_UpdateNoteByID(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	note := &models.Note{
		ID:      1,
		Title:   "First Note",
		Content: "This is the first note",
	}

	repo := NewNotesRepository(db)

	mock.ExpectExec("UPDATE notes").WithArgs(note.Title, note.Content, note.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdateNoteByID(note.ID, note)

	// Assert
	assert.NoError(t, err)
}

func TestNoteRepository_DeleteNoteByID(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	note := &models.Note{
		ID:      1,
		Title:   "First Note",
		Content: "This is the first note",
	}

	repo := NewNotesRepository(db)

	mock.ExpectExec("DELETE FROM notes").WithArgs(note.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.DeleteNoteByID(note.ID)

	// Assert
	assert.NoError(t, err)
}

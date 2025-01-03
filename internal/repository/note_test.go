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
	firstId, firstErr := repo.Create(firstNote)
	secondId, secondErr := repo.Create(secondNote)

	// Assert
	assert.NoError(t, firstErr)
	assert.NoError(t, secondErr)
	assert.Equal(t, 1, firstId)
	assert.Equal(t, 2, secondId)

}

func TestNoteRepository_GetNoteById(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening a stub database connection: %s", err)
	}
	defer db.Close()
	notes := [2]*models.Note{
		&models.Note{
			Id:      1,
			Title:   "First Note",
			Content: "This is the first note",
		},
		&models.Note{
			Id:      2,
			Title:   "Second Note",
			Content: "This is the second note",
		}}

	repo := NewNotesRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "content"})
	for _, note := range notes {
		rows.AddRow(note.Id, note.Title, note.Content)
	}

	mock.ExpectQuery("SELECT id, title, content FROM notes WHERE id = ?").WithArgs(notes[0].Id).WillReturnRows(rows)
	mock.ExpectQuery("SELECT id, title, content FROM notes WHERE id = ?").WithArgs(notes[1].Id).WillReturnRows(rows)

	// Act
	firstResult, firstErr := repo.Get(notes[0].Id)
	secondResult, secondErr := repo.Get(notes[1].Id)

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
			Id:      1,
			Title:   "First Note",
			Content: "This is the first note",
		},
		&models.Note{
			Id:      2,
			Title:   "Second Note",
			Content: "This is the second note",
		},
		&models.Note{
			Id:      3,
			Title:   "Third Note",
			Content: "This is the third note",
		},
	}

	repo := NewNotesRepository(db)

	rows := sqlmock.NewRows([]string{"id", "title", "content"})
	for _, note := range notes {
		rows.AddRow(note.Id, note.Title, note.Content)
	}

	mock.ExpectQuery("SELECT id, title, content FROM notes").WillReturnRows(rows)

	// Act
	res, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, len(res), 3)
	for i := range notes {
		assert.Equal(t, res[i], notes[i])
	}
}

func TestNoteRepository_UpdateNoteById(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	note := &models.Note{
		Id:      1,
		Title:   "First Note",
		Content: "This is the first note",
	}

	repo := NewNotesRepository(db)

	mock.ExpectExec("UPDATE notes").WithArgs(note.Title, note.Content, note.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.Update(note.Id, note)

	// Assert
	assert.NoError(t, err)
}

func TestNoteRepository_DeleteNoteById(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	note := &models.Note{
		Id:      1,
		Title:   "First Note",
		Content: "This is the first note",
	}

	repo := NewNotesRepository(db)

	mock.ExpectExec("DELETE FROM notes").WithArgs(note.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.Delete(note.Id)

	// Assert
	assert.NoError(t, err)
}

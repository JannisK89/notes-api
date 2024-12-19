package mocks

import (
	"github.com/JannisK89/notes-api/internal/models"
	"github.com/stretchr/testify/mock"
)

type NoteRepoMock struct {
	mock.Mock
}

func (m *NoteRepoMock) CreateNote(note *models.Note) (int, error) {
	args := m.Called(note)
	return args.Int(0), args.Error(1)
}

func (m *NoteRepoMock) GetNoteByID(id int) (*models.Note, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Note), args.Error(1)
}

func (m *NoteRepoMock) GetAllNotes() ([]*models.Note, error) {
	args := m.Called()
	return args.Get(0).([]*models.Note), args.Error(1)
}

func (m *NoteRepoMock) UpdateNoteByID(id int, note *models.Note) error {
	args := m.Called(id, note)
	return args.Error(0)
}

func (m *NoteRepoMock) DeleteNoteByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

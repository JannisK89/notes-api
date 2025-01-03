package mocks

import (
	"github.com/JannisK89/notes-api/internal/models"
	"github.com/stretchr/testify/mock"
)

// NoteRepoMock is a mock for the NoteRepository interface
type NoteRepoMock struct {
	mock.Mock
}

// Get mocks the Get method of the NoteRepository interface
func (m *NoteRepoMock) Get(id int) (*models.Note, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Note), args.Error(1)
}

// GetAll mocks the GetAll method of the NoteRepository interface
func (m *NoteRepoMock) GetAll() ([]*models.Note, error) {
	args := m.Called()
	return args.Get(0).([]*models.Note), args.Error(1)
}

// Create mocks the Create method of the NoteRepository interface
func (m *NoteRepoMock) Create(note *models.Note) (int, error) {
	args := m.Called(note)
	return args.Int(0), args.Error(1)
}

// Update mocks the Update method of the NoteRepository interface
func (m *NoteRepoMock) Update(id int, note *models.Note) error {
	args := m.Called(id, note)
	return args.Error(0)
}

// Delete mocks the Delete method of the NoteRepository interface
func (m *NoteRepoMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

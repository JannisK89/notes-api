package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JannisK89/notes-api/internal/mocks"
	"github.com/JannisK89/notes-api/internal/models"
	"github.com/JannisK89/notes-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// TODO: Fix Service Mock
// TODO: Add failure tests
func TestNoteHandler_Get(t *testing.T) {
	// Arrange
	noteRepoMock := &mocks.NoteRepoMock{}
	noteService := service.NewNoteService(noteRepoMock)
	noteHandler := NewNoteHandler(noteService)

	note := &models.Note{Id: 1, Title: "Test Note", Content: "I Am A Test Note"}

	noteRepoMock.On("Get", 1).Return(note, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notes/", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("noteId", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	noteHandler.Get(rec, req)

	// Assertion
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status": "ok", "data": {"id":1,"title":"Test Note","content":"I Am A Test Note"} }`, rec.Body.String())
	noteRepoMock.AssertExpectations(t)
}

func TestNoteHandler_Create(t *testing.T) {
	// Arrange
	noteRepoMock := &mocks.NoteRepoMock{}
	noteService := service.NewNoteService(noteRepoMock)
	noteHandler := NewNoteHandler(noteService)

	note := &models.Note{Title: "Test Note", Content: "I Am A Test Note"}
	noteRepoMock.On("Create", note).Return(1, nil)

	payload, err := json.Marshal(note)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	// Act
	noteHandler.Create(rec, req)

	// Assertion
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"status": "ok", "data": 1}`, rec.Body.String())
	noteRepoMock.AssertExpectations(t)
}

func TestNoteHandler_Update(t *testing.T) {
	// Arrange
	noteRepoMock := &mocks.NoteRepoMock{}
	noteService := service.NewNoteService(noteRepoMock)
	noteHandler := NewNoteHandler(noteService)

	note := &models.Note{Id: 1, Title: "Test Note", Content: "I Am A Test Note"}
	noteRepoMock.On("Update", 1, note).Return(nil)

	payload, err := json.Marshal(note)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/notes", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("noteId", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	noteHandler.Update(rec, req)

	// Assertion
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status": "ok", "message": "Note Updated"}`, rec.Body.String())
	noteRepoMock.AssertExpectations(t)
}

func TestNoteHandler_Delete(t *testing.T) {
	// Arrange
	noteRepoMock := &mocks.NoteRepoMock{}
	noteService := service.NewNoteService(noteRepoMock)
	noteHandler := NewNoteHandler(noteService)

	noteRepoMock.On("Delete", 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notes", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("noteId", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	noteHandler.Delete(rec, req)

	// Assertion
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status": "ok", "message": "Note Deleted"}`, rec.Body.String())
	noteRepoMock.AssertExpectations(t)
}

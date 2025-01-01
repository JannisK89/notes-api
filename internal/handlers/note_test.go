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
func TestNoteHandler_GetNote_Success(t *testing.T) {
	// Arrange
	noteRepoMock := &mocks.NoteRepoMock{}
	noteService := service.NewNoteService(noteRepoMock)
	noteHandler := NewNoteHandler(noteService)

	note := &models.Note{ID: 1, Title: "Test Note", Content: "I Am A Test Note"}

	noteRepoMock.On("GetNoteByID", 1).Return(note, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notes/", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("noteID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	noteHandler.GetNote(rec, req)

	// Assertion
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status": "success", "data": {"id":1,"title":"Test Note","content":"I Am A Test Note"} }`, rec.Body.String())
	noteRepoMock.AssertExpectations(t)
}

func TestNoteHandler_CreateNote_Success(t *testing.T) {
	// Arrange
	noteRepoMock := &mocks.NoteRepoMock{}
	noteService := service.NewNoteService(noteRepoMock)
	noteHandler := NewNoteHandler(noteService)

	note := &models.Note{Title: "Test Note", Content: "I Am A Test Note"}
	noteRepoMock.On("CreateNote", note).Return(1, nil)

	payload, err := json.Marshal(note)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	// Act
	noteHandler.CreateNote(rec, req)

	// Assertion
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"status": "success", "data": 1}`, rec.Body.String())
	noteRepoMock.AssertExpectations(t)
}

func TestNoteHandler_UpdateNote_Success(t *testing.T) {
	// Arrange
	noteRepoMock := &mocks.NoteRepoMock{}
	noteService := service.NewNoteService(noteRepoMock)
	noteHandler := NewNoteHandler(noteService)

	note := &models.Note{ID: 1, Title: "Test Note", Content: "I Am A Test Note"}
	noteRepoMock.On("UpdateNoteByID", 1, note).Return(nil)

	payload, err := json.Marshal(note)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/notes", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("noteID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	noteHandler.UpdateNote(rec, req)

	// Assertion
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status": "success", "message": "Note Updated"}`, rec.Body.String())
	noteRepoMock.AssertExpectations(t)
}

func TestNoteHandler_DeleteNote_Success(t *testing.T) {
	// Arrange
	noteRepoMock := &mocks.NoteRepoMock{}
	noteService := service.NewNoteService(noteRepoMock)
	noteHandler := NewNoteHandler(noteService)

	noteRepoMock.On("DeleteNoteByID", 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notes", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("noteID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	noteHandler.DeleteNote(rec, req)

	// Assertion
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.JSONEq(t, `{"status": "success"}`, rec.Body.String())
	noteRepoMock.AssertExpectations(t)
}

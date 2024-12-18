package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JannisK89/notes-api/internal/models"
	"github.com/JannisK89/notes-api/internal/repository"
	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	repo repository.NoteRepository
}

func NewNoteHandler(repo repository.NoteRepository) *NoteHandler {
	return &NoteHandler{repo}
}

func (h NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteID")
	noteIDAsInt, err := strconv.Atoi(noteID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("noteID must be an integer"))
		return
	}

	note, err := h.repo.GetNoteByID(noteIDAsInt)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Note not found"))
		return
	}

	jsonNote, err := json.Marshal(note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(jsonNote)
}

func (h NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	note := &models.Note{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return
	}

	id, err := h.repo.CreateNote(note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("{\"id\":" + strconv.Itoa(id) + "}"))
}

func (h NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.repo.GetAllNotes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	jsonNotes, err := json.Marshal(notes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(jsonNotes)
}

func (h NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteID")
	noteIDAsInt, err := strconv.Atoi(noteID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("noteID must be an integer"))
		return
	}

	note := &models.Note{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return
	}

	err = h.repo.UpdateNoteByID(noteIDAsInt, note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Note updated"))
}

func (h NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteID")
	noteIDAsInt, err := strconv.Atoi(noteID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("noteID must be an integer"))
		return
	}

	err = h.repo.DeleteNoteByID(noteIDAsInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

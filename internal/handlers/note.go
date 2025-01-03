package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/JannisK89/notes-api/internal/models"
	"github.com/JannisK89/notes-api/internal/repository"
	"github.com/JannisK89/notes-api/internal/service"
	"github.com/JannisK89/notes-api/internal/utils"
	"github.com/go-chi/chi/v5"
)

// ErrInvalidId is returned when the id is not a valid integer
var ErrInvalidId = errors.New("id must be a valid integer")

// getNoteid extracts the noteid from the URL and returns it as an integer
// It returns an error if the noteid is not a valid integer
func getNoteId(r *http.Request) (int, error) {
	noteid := chi.URLParam(r, "noteId")
	noteidAsInt, err := strconv.Atoi(noteid)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidId, noteid)
	}
	return noteidAsInt, nil
}

// NoteHandler handles HTTP requests related to notes. It provides methods for
// creating, retrieving, updating, and deleting notes.
type NoteHandler struct {
	noteService service.NoteService
}

// NewNoteHandler creates a new NoteHandler
func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{noteService}
}

// Get retrieves a note by its id from the database.
// It returns a 404 error if the note is not found and a 400 error
// if the provided id is not a valid integer.
func (h NoteHandler) Get(w http.ResponseWriter, r *http.Request) {
	noteid, err := getNoteId(r)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, ErrInvalidId.Error())
		return
	}

	note, err := h.noteService.Get(noteid)
	if err != nil {
		log.Println(err)
		if errors.Is(err, repository.ErrNoteNotFound) {
			utils.ErrorResponse(w, http.StatusNotFound, repository.ErrNoteNotFound.Error())
			return
		} else if errors.Is(err, service.ErrInvalidId) {
			utils.ErrorResponse(w, http.StatusBadRequest, service.ErrInvalidId.Error())
			return

		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, utils.InternalServerError)
			return
		}
	}
	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Status: utils.StatusOk, Data: note})
}

// Create adds a new note to the database.
// It returns a 400 error if the note data is invalid.
func (h NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	note := &models.Note{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, utils.BadRequest)
		return
	}

	id, err := h.noteService.Create(note)
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrInvalidNote) {
			utils.ErrorResponse(w, http.StatusBadRequest, service.ErrInvalidNote.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, utils.InternalServerError)
		return

	}
	utils.JSONResponse(w, http.StatusCreated, utils.ApiResponse{Status: utils.StatusOk, Data: id})
}

// GetAll retrieves all notes from the database.
func (h NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteService.GetAll()
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Status: utils.StatusOk, Data: notes})
}

// Update modifies an existing note in the database.
// It returns a 400 error if the note data is invalid or if the id is not found.
func (h NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	noteid, err := getNoteId(r)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	note := &models.Note{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, utils.BadRequest)
		return
	}

	err = h.noteService.Update(noteid, note)
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrInvalidNote) {
			utils.ErrorResponse(w, http.StatusBadRequest, service.ErrInvalidNote.Error())
			return
		} else if errors.Is(err, service.ErrInvalidId) {
			utils.ErrorResponse(w, http.StatusBadRequest, service.ErrInvalidId.Error())
			return
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, utils.InternalServerError)
			return
		}
	}

	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Message: "Note Updated", Status: utils.StatusOk})
}

// Delete removes a note from the database.
// It returns a 400 error if the id is not found.
func (h NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	noteid, err := getNoteId(r)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.noteService.Delete(noteid)
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrInvalidId) {
			utils.ErrorResponse(w, http.StatusBadRequest, service.ErrInvalidId.Error())
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return

	}
	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Status: utils.StatusOk, Message: "Note Deleted"})
}

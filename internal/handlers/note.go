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

var ErrInvalidID = errors.New("ID must be a valid integer")

func getNoteId(r *http.Request) (int, error) {
	noteId := chi.URLParam(r, "noteId")
	noteIdAsInt, err := strconv.Atoi(noteId)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidID, noteId)
	}
	return noteIdAsInt, nil
}

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{noteService}
}

func (h NoteHandler) Get(w http.ResponseWriter, r *http.Request) {
	noteId, err := getNoteId(r)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, ErrInvalidID.Error())
		return
	}

	note, err := h.noteService.Get(noteId)
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
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, utils.InternalServerError)
			return
		}
	}
	utils.JSONResponse(w, http.StatusCreated, utils.ApiResponse{Status: utils.StatusOk, Data: id})
}

func (h NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteService.GetAll()
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Status: utils.StatusOk, Data: notes})
}

func (h NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	noteId, err := getNoteId(r)
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

	err = h.noteService.Update(noteId, note)
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

func (h NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	noteId, err := getNoteId(r)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.noteService.Delete(noteId)
	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrInvalidId) {
			utils.ErrorResponse(w, http.StatusBadRequest, service.ErrInvalidId.Error())
			return
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}
	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Status: utils.StatusOk, Message: "Note Deleted"})
}

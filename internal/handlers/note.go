package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/JannisK89/notes-api/internal/models"
	"github.com/JannisK89/notes-api/internal/repository"
	"github.com/JannisK89/notes-api/internal/service"
	"github.com/JannisK89/notes-api/internal/utils"
	"github.com/go-chi/chi/v5"
)

var (
	statusSuccess = "success"
	statusError   = "error"
)

var ErrInvalidID = errors.New("ID must be a valid integer")

func getNoteId(r *http.Request) (int, error) {
	noteId := chi.URLParam(r, "noteID")
	noteIdAsInt, err := strconv.Atoi(noteId)
	if err != nil {
		return 0, ErrInvalidID
	}
	return noteIdAsInt, nil
}

type NoteHandler struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) *NoteHandler {
	return &NoteHandler{noteService}
}

func (h NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	noteId, err := getNoteId(r)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	note, err := h.noteService.GetNote(noteId)
	if err != nil {
		log.Println(err)
		if repoError, ok := err.(*repository.RepoError); ok && repoError.Err == repository.ErrNoteNotFound {
			utils.ErrorResponse(w, http.StatusNotFound, repoError.Err.Error())
			return
		} else if serviceError, ok := err.(*service.ServiceError); ok {
			utils.ErrorResponse(w, http.StatusBadRequest, serviceError.Err.Error())
			return

		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}
	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Status: statusSuccess, Data: note})
}

func (h NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	note := &models.Note{}
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	id, err := h.noteService.CreateNote(note)
	if err != nil {
		log.Println(err)
		if serviceError, ok := err.(*service.ServiceError); ok {
			utils.ErrorResponse(w, http.StatusBadRequest, serviceError.Err.Error())
			return
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}
	utils.JSONResponse(w, http.StatusCreated, utils.ApiResponse{Status: statusSuccess, Data: id})
}

func (h NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteService.GetAllNotes()
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Status: statusSuccess, Data: notes})
}

func (h NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
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
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid Data")
		return
	}

	err = h.noteService.UpdateNote(noteId, note)
	if err != nil {
		log.Println(err)
		if serviceError, ok := err.(*service.ServiceError); ok {
			utils.ErrorResponse(w, http.StatusBadRequest, serviceError.Err.Error())
			return
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Message: "Note Updated", Status: statusSuccess})
}

func (h NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteId, err := getNoteId(r)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.noteService.DeleteNote(noteId)
	if err != nil {
		log.Println(err)
		if serviceError, ok := err.(*service.ServiceError); ok {
			utils.ErrorResponse(w, http.StatusBadRequest, serviceError.Err.Error())
			return
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}
	utils.JSONResponse(w, http.StatusNoContent, utils.ApiResponse{Status: statusSuccess})

}

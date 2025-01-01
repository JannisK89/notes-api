package handlers

import (
	"encoding/json"
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

var (
	statusSuccess = "success"
	statusError   = "error"
)

func getNoteId(r *http.Request) (int, error) {
	noteId := chi.URLParam(r, "noteID")
	noteIdAsInt, err := strconv.Atoi(noteId)
	if err != nil {
		return 0, fmt.Errorf("noteId must be a valid integer: %v", noteId)
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
		utils.JSONResponse(w, http.StatusBadRequest, utils.ApiResponse{Message: err.Error(), Status: statusError})
		return
	}

	note, err := h.noteService.GetNote(noteId)
	if err != nil {
		log.Println(err)
		if repoError, ok := err.(*repository.RepoError); ok && repoError.Err == repository.ErrNoteNotFound {
			utils.JSONResponse(w, http.StatusNotFound, utils.ApiResponse{Message: repoError.Err.Error(), Status: statusError})
			return
		} else {
			utils.JSONResponse(w, http.StatusInternalServerError, utils.ApiResponse{Message: "Internal Server Error", Status: statusError})
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
		utils.JSONResponse(w, http.StatusBadRequest, utils.ApiResponse{Message: "Invalid Request", Status: statusError})
		return
	}

	id, err := h.noteService.CreateNote(note)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, utils.ApiResponse{Message: "Internal Server Error", Status: statusError})
		return
	}
	utils.JSONResponse(w, http.StatusCreated, utils.ApiResponse{Status: statusSuccess, Data: id})
}

func (h NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.noteService.GetAllNotes()
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, utils.ApiResponse{Message: "Internal Server Error", Status: statusError})
		return
	}

	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Status: statusSuccess, Data: notes})
}

func (h NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteID")
	noteIDAsInt, err := strconv.Atoi(noteID)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusBadRequest, utils.ApiResponse{Message: err.Error(), Status: statusError})
		return
	}

	note := &models.Note{}
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(note)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusBadRequest, utils.ApiResponse{Status: statusError, Message: "Invalid Data"})
		return
	}

	err = h.noteService.UpdateNote(noteIDAsInt, note)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, utils.ApiResponse{Message: "Internal Server Error", Status: statusError})
		return
	}

	utils.JSONResponse(w, http.StatusOK, utils.ApiResponse{Message: "Note Updated", Status: statusSuccess})
}

func (h NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteID")
	noteIDAsInt, err := strconv.Atoi(noteID)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusBadRequest, utils.ApiResponse{Message: err.Error(), Status: statusError})
		return
	}

	err = h.noteService.DeleteNote(noteIDAsInt)
	if err != nil {
		log.Println(err)
		utils.JSONResponse(w, http.StatusInternalServerError, utils.ApiResponse{Message: "Internal Server Error", Status: statusError})
		return
	}

	utils.JSONResponse(w, http.StatusNoContent, utils.ApiResponse{Status: statusSuccess})

}

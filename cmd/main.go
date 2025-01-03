package main

import (
	"log"
	"net/http"

	"github.com/JannisK89/notes-api/internal/db"
	"github.com/JannisK89/notes-api/internal/handlers"
	"github.com/JannisK89/notes-api/internal/repository"
	"github.com/JannisK89/notes-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	dbconn, err := db.NewSQLiteDB("./notes.db")
	if err != nil {
		log.Fatal("Could not connect to Database: ", err)
	}

	notesRepo := repository.NewNotesRepository(dbconn)
	notesService := service.NewNoteService(notesRepo)
	notesHandler := handlers.NewNoteHandler(notesService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/notes", func(r chi.Router) {
			r.Get("/", notesHandler.GetAll)
			r.Post("/", notesHandler.Create)
			r.Get("/{noteId}", notesHandler.Get)
			r.Put("/{noteId}", notesHandler.Update)
			r.Delete("/{noteId}", notesHandler.Delete)
		})
	})

	error := http.ListenAndServe(":3000", r)
	if error != nil {
		log.Fatal("Could not start server: ", error)
	}

}

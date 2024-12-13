package main

import (
	"log"
	"net/http"

	"github.com/JannisK89/notes-api/internal/db"
	"github.com/JannisK89/notes-api/internal/handlers"
	"github.com/JannisK89/notes-api/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	dbconn, err := db.NewSQLiteDB("./notes.db")
	if err != nil {
		log.Fatal("Could not connect to Database: ", err)
	}

	notesRepo := repository.NewNotesRepository(dbconn)
	notesHandler := handlers.NewNoteHandler(notesRepo)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/notes", func(r chi.Router) {
			r.Get("/", notesHandler.GetAllNotes)
			r.Get("/{noteID}", notesHandler.GetNote)
			r.Post("/", notesHandler.CreateNote)
		})
	})

	http.ListenAndServe(":3000", r)
}

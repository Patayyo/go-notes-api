package main

import (
	"log"
	"net/http"
	"notes-api/db"
	"notes-api/handler"
	storage "notes-api/repo"
	"notes-api/service"

	"github.com/gorilla/mux"
)

func main() {
	db.ConnectDB()
	newStore := storage.NewPostgresStore(db.DB)
	service := service.NewNoteService(newStore)
	h := &handler.NoteHandler{Store: service}

	r := mux.NewRouter()
	r.HandleFunc("/notes", h.GetAll).Methods("GET")
	r.HandleFunc("/notes/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/notes", h.Create).Methods("POST")
	r.HandleFunc("/notes/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/notes/{id}", h.Delete).Methods("DELETE")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

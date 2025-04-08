package main

import (
	"log"
	"net/http"
	"notes-api/auth"
	"notes-api/db"
	"notes-api/handler"
	"notes-api/middleware"
	storage "notes-api/repo"
	"notes-api/service"

	"github.com/gorilla/mux"
)

func main() {
	db.ConnectDB()
	newStore := storage.NewPostgresStore(db.DB)
	service := service.NewNoteService(newStore)
	authService := auth.NewAuthService(db.DB)
	h := &handler.NoteHandler{Store: service}
	authHandler := &auth.AuthHandler{Service: authService}

	r := mux.NewRouter()

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	authRoutes := r.PathPrefix("/notes").Subrouter()
	authRoutes.Use(middleware.JWTAuthMiddleware)
	authRoutes.HandleFunc("", h.GetAll).Methods("GET")
	authRoutes.HandleFunc("", h.Create).Methods("POST")
	authRoutes.HandleFunc("/{id}", h.GetByID).Methods("GET")
	authRoutes.HandleFunc("/{id}", h.Update).Methods("PUT")
	authRoutes.HandleFunc("/{id}", h.Delete).Methods("DELETE")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

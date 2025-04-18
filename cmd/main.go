// @title Notes API
// @version 1.0
// @description API для управления заметками пользователя.
// @host localhost:8080
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"net/http"
	"notes-api/auth"
	"notes-api/db"
	"notes-api/handler"
	"notes-api/logger"
	"notes-api/middleware"
	storage "notes-api/repo"
	"notes-api/service"

	_ "notes-api/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	logger.Init()
	logger.Log.Info("Логгер инициализирован")

	db.ConnectDB()
	newStore := storage.NewPostgresStore(db.DB)
	service := service.NewNoteService(newStore)
	authService := auth.NewAuthService(db.DB)
	h := &handler.NoteHandler{Store: service}
	authHandler := &auth.AuthHandler{Service: authService}

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")

	authRoutes := r.PathPrefix("/notes").Subrouter()
	authRoutes.Use(middleware.JWTAuthMiddleware)
	authRoutes.HandleFunc("", h.GetAll).Methods("GET")
	authRoutes.HandleFunc("", h.Create).Methods("POST")
	authRoutes.HandleFunc("/{id}", h.GetByID).Methods("GET")
	authRoutes.HandleFunc("/{id}", h.Update).Methods("PUT")
	authRoutes.HandleFunc("/{id}", h.Delete).Methods("DELETE")

	authProtected := r.NewRoute().Subrouter()
	authProtected.Use(middleware.JWTAuthMiddleware)
	authProtected.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

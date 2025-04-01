package service

import (
	"errors"
	"notes-api/model"
	storage "notes-api/repo"
)

type NoteService struct {
	Repo storage.NoteRepository
}

type INoteService interface {
	GetAllNotes() ([]model.Note, error)
	GetNoteByID(id int) (model.Note, error)
	CreateNote(note model.Note) (model.Note, error)
	UpdateNote(id int, updated model.Note) (model.Note, error)
	DeleteNote(id int) error
}

func NewNoteService(r storage.NoteRepository) *NoteService {
	return &NoteService{Repo: r}
}

func (s *NoteService) GetAllNotes() ([]model.Note, error) {
	return s.Repo.GetAll()
}

func (s *NoteService) GetNoteByID(id int) (model.Note, error) {
	note, err := s.Repo.GetByID(id)
	if err != nil {
		return model.Note{}, err
	}
	return note, nil
}

func (s *NoteService) CreateNote(note model.Note) (model.Note, error) {
	if note.Title == "" && note.Content == "" {
		return model.Note{}, errors.New("заголовок и содержание не могут быть пустыми")
	}
	return s.Repo.Create(note)
}

func (s *NoteService) UpdateNote(id int, updated model.Note) (model.Note, error) {
	if updated.Title == "" && updated.Content == "" {
		return model.Note{}, errors.New("заголовок и содержание не могут быть пустыми")
	}
	return s.Repo.Update(id, updated)
}

func (s *NoteService) DeleteNote(id int) error {
	return s.Repo.Delete(id)
}

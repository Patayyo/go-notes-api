package storage

import (
	"notes-api/model"

	"gorm.io/gorm"
)

type PostgresStore struct {
	DB *gorm.DB
}

type NoteRepository interface {
	GetAll() ([]model.Note, error)
	GetByID(id int) (model.Note, error)
	Create(note model.Note) (model.Note, error)
	Update(id int, updated model.Note) (model.Note, error)
	Delete(id int) error
}

func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{DB: db}
}

func (s *PostgresStore) GetAll() ([]model.Note, error) {
	var notes []model.Note
	if err := s.DB.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *PostgresStore) GetByID(id int) (model.Note, error) {
	var note model.Note
	if err := s.DB.First(&note, id).Error; err != nil {
		return model.Note{}, err
	}
	return note, nil
}

func (s *PostgresStore) Create(note model.Note) (model.Note, error) {
	if err := s.DB.Create(&note).Error; err != nil {
		return model.Note{}, err
	}
	return note, nil
}

func (s *PostgresStore) Update(id int, updated model.Note) (model.Note, error) {
	var note model.Note
	if err := s.DB.Find(&note, id).Error; err != nil {
		return model.Note{}, err
	}
	note.Title = updated.Title
	note.Content = updated.Content
	if err := s.DB.Save(&note).Error; err != nil {
		return model.Note{}, err
	}
	return note, nil
}

func (s *PostgresStore) Delete(id int) error {
	var note model.Note
	if err := s.DB.Find(&note, id).Error; err != nil {
		return err
	}
	if err := s.DB.Delete(&note).Error; err != nil {
		return err
	}
	return nil
}

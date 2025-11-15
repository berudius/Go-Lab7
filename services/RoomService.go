package services

import (
	"go.mod/models"
	"go.mod/repositories"
)

type RoomService interface {
	GetAll() ([]models.Room, error)
	GetByID(id uint) (models.Room, error)
	Create(room *models.Room) error
	Update(room *models.Room) error
	Delete(id uint) error
}

type roomServiceImpl struct {
	repo repositories.RoomRepository
}

func NewRoomService(repo repositories.RoomRepository) RoomService {
	return &roomServiceImpl{repo: repo}
}

func (s *roomServiceImpl) GetAll() ([]models.Room, error) {
	return s.repo.GetAll()
}

func (s *roomServiceImpl) GetByID(id uint) (models.Room, error) {
	return s.repo.GetByID(id)
}

func (s *roomServiceImpl) Create(room *models.Room) error {
	return s.repo.Create(room)
}

func (s *roomServiceImpl) Update(room *models.Room) error {
	return s.repo.Update(room)
}

func (s *roomServiceImpl) Delete(id uint) error {
	return s.repo.Delete(id)
}

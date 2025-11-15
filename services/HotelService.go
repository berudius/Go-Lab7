package services

import (
	"go.mod/models"
	"go.mod/repositories"
)

type HotelService interface {
	GetAll() ([]models.Hotel, error)
	GetByID(id uint) (models.Hotel, error)
	Create(hotel *models.Hotel) error
	Update(hotel *models.Hotel) error
	Delete(id uint) error
}

type hotelServiceImpl struct {
	repo repositories.HotelRepository
}

func NewHotelService(repo repositories.HotelRepository) HotelService {
	return &hotelServiceImpl{repo: repo}
}

func (s *hotelServiceImpl) GetAll() ([]models.Hotel, error) {
	return s.repo.GetAll()
}

func (s *hotelServiceImpl) GetByID(id uint) (models.Hotel, error) {
	return s.repo.GetByID(id)
}

func (s *hotelServiceImpl) Create(hotel *models.Hotel) error {
	return s.repo.Create(hotel)
}

func (s *hotelServiceImpl) Update(hotel *models.Hotel) error {
	return s.repo.Update(hotel)
}

func (s *hotelServiceImpl) Delete(id uint) error {
	return s.repo.Delete(id)
}
